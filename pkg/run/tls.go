package run

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"
)

func genServerCert() (caCert, serverCert *tls.Certificate, err error) {
	// CA cert
	caClaims := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "local-test",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(0, 1, 0),
		IsCA:      true,
	}
	caKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, fmt.Errorf("ca keygen: %w", err)
	}
	caCertBytes, err := x509.CreateCertificate(rand.Reader, caClaims, caClaims, caKey.Public(), caKey)
	if err != nil {
		return nil, nil, fmt.Errorf("ca certgen: %w", err)
	}
	caCertPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caCertBytes,
	})
	caKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caKey),
	})
	caKeyPair, err := tls.X509KeyPair(caCertPem, caKeyPem)
	if err != nil {
		return nil, nil, fmt.Errorf("ca keypair: %w", err)
	}
	caCert = &caKeyPair

	// server cert
	serverClaims := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		DNSNames:    []string{"localhost", "localhost.localdomain", "[::1]"},
		NotBefore:   caClaims.NotBefore,
		NotAfter:    caClaims.NotAfter,
	}
	serverKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, fmt.Errorf("server keygen: %w", err)
	}
	serverCertBytes, err := x509.CreateCertificate(rand.Reader, serverClaims, caClaims, serverKey.Public(), caKey)
	if err != nil {
		return nil, nil, fmt.Errorf("server cert: %w", err)
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
		return nil, nil, fmt.Errorf("server keypair: %w", err)
	}
	serverCert = &serverKeyPair

	return caCert, serverCert, nil
}
