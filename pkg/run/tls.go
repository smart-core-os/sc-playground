package run

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"time"
)

// ca manages and generates certificates for use by TLS.
type ca struct {
	key        *rsa.PrivateKey
	cert       *x509.Certificate
	x509Cert   []byte          // contains cert, signed by key
	certPem    []byte          // encoded x509Cert
	keyPem     []byte          // encoded key
	tlsKeyPair tls.Certificate // combines certPem and keyPem

	serial      *big.Int // tracks certs signed by key, server and client certs
	clientCount int      // tracks how many clients we generated, so we can give them unique names
}

func NewSelfSignedCA() (*ca, error) {
	// CA cert
	ca := &ca{
		serial: big.NewInt(0),
	}
	var err error
	ca.cert = &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "sc-playground-ca",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(0, 1, 0),

		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}
	ca.key, err = rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, fmt.Errorf("keygen: %w", err)
	}
	ca.x509Cert, err = x509.CreateCertificate(rand.Reader, ca.cert, ca.cert, ca.key.Public(), ca.key)
	if err != nil {
		return nil, fmt.Errorf("sign: %w", err)
	}
	ca.certPem = pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: ca.x509Cert,
	})
	ca.keyPem = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(ca.key),
	})
	ca.tlsKeyPair, err = tls.X509KeyPair(ca.certPem, ca.keyPem)
	if err != nil {
		return nil, fmt.Errorf("ca keypair: %w", err)
	}
	return ca, nil
}

func (ca *ca) nextSerial() *big.Int {
	return ca.serial.Add(ca.serial, big.NewInt(1))
}

func (ca *ca) Pool() (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(ca.certPem) {
		return nil, errors.New("failed to append cert")
	}
	return pool, nil
}

func (ca *ca) WriteCACertPEM(w io.Writer) error {
	_, err := w.Write(ca.certPem)
	return err
}

func (ca *ca) NewServerCert() (*tls.Certificate, error) {
	// server cert
	serverClaims := &x509.Certificate{
		SerialNumber: ca.nextSerial(),
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		DNSNames:    []string{"localhost", "localhost.localdomain", "[::1]"},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(0, 1, 0),

		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	serverKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, fmt.Errorf("keygen: %w", err)
	}
	serverCertBytes, err := x509.CreateCertificate(rand.Reader, serverClaims, ca.cert, serverKey.Public(), ca.key)
	if err != nil {
		return nil, fmt.Errorf("sign: %w", err)
	}
	serverCertPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: serverCertBytes,
	})
	serverKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(serverKey),
	})

	serverKeyPair, err := tls.X509KeyPair(serverCertPem, serverKeyPem)
	if err != nil {
		return nil, fmt.Errorf("keypair: %w", err)
	}
	return &serverKeyPair, nil
}

func (ca *ca) WriteClientCert(w io.Writer) error {
	// client cert
	ca.clientCount++
	name := fmt.Sprintf("client-%d", ca.clientCount)
	clientClaims := &x509.Certificate{
		SerialNumber: ca.nextSerial(),
		Subject: pkix.Name{
			CommonName: name,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(0, 1, 0),

		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	clientKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("keygen: %w", err)
	}
	clientCertBytes, err := x509.CreateCertificate(rand.Reader, clientClaims, ca.cert, clientKey.Public(), ca.key)
	if err != nil {
		return fmt.Errorf("sign: %w", err)
	}
	if err := pem.Encode(w, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: clientCertBytes,
	}); err != nil {
		return fmt.Errorf("write cert: %w", err)
	}
	if err := pem.Encode(w, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(clientKey),
	}); err != nil {
		return fmt.Errorf("write key: %w", err)
	}

	return nil
}
