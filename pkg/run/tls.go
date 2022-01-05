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
	"os"
	"path/filepath"
	"time"
)

const (
	caCertFile     = "ca-cer.pem"
	caKeyFile      = "ca-key.pem"
	serverCertFile = "server%d-cer.pem"
	serverKeyFile  = "server%d-key.pem"
)

// ca manages and generates certificates for use by TLS.
type ca struct {
	key        *rsa.PrivateKey
	cert       *x509.Certificate
	certPem    []byte          // encoded x509Cert
	keyPem     []byte          // encoded key
	tlsKeyPair tls.Certificate // combines certPem and keyPem

	serial      *big.Int // tracks certs signed by key, server and client certs
	clientCount int      // tracks how many clients we generated, so we can give them unique names
	serverCount int      // tracks how many servers we generated, so we can give them unique names and cache correctly

	cacheDir        string // for caching/loading certs
	cacheReadFailed bool   // set to true if reading the CA from the cache failed for this instance.
}

func LoadOrCreateSelfSignedCA(opts ...caOption) (instance *ca, fromCache bool, err error) {
	// CA cert
	ca := &ca{}
	for _, opt := range opts {
		opt(ca)
	}

	if ca.cacheDir != "" {
		if err := ca.loadCAFromCache(); err != nil {
			// if we error loading from the cache, continue with generating a new one
			if !errors.Is(err, os.ErrNotExist) {
				fmt.Printf("Cache load: %v\n", err.Error())
			}
			ca.cacheReadFailed = true
		} else {
			return ca, true, nil
		}
	}

	ca.serial = big.NewInt(0)
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
		return nil, false, fmt.Errorf("keygen: %w", err)
	}
	x509Cert, err := x509.CreateCertificate(rand.Reader, ca.cert, ca.cert, ca.key.Public(), ca.key)
	if err != nil {
		return nil, false, fmt.Errorf("sign: %w", err)
	}
	ca.certPem = pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: x509Cert,
	})
	ca.keyPem = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(ca.key),
	})
	ca.tlsKeyPair, err = tls.X509KeyPair(ca.certPem, ca.keyPem)
	if err != nil {
		return nil, false, fmt.Errorf("ca keypair: %w", err)
	}

	if err := ca.writeCAToCache(); err != nil {
		fmt.Printf("Error writing CA cache: %v", err)
	}
	return ca, false, nil
}

func (ca *ca) loadCAFromCache() error {
	// attempt to load the ca
	caKeyPem, err := os.ReadFile(filepath.Join(ca.cacheDir, caKeyFile))
	if err != nil {
		return fmt.Errorf("key %w", err)
	}
	caCertPem, err := os.ReadFile(filepath.Join(ca.cacheDir, caCertFile))
	if err != nil {
		return fmt.Errorf("cert %w", err)
	}
	x509KeyPair, err := tls.X509KeyPair(caCertPem, caKeyPem)
	if err != nil {
		return fmt.Errorf("keypair parse %w", err)
	}

	// if len(x509KeyPair.Certificate) == 0 then tls.X509KeyPair returns an error
	cert, err := x509.ParseCertificate(x509KeyPair.Certificate[0])
	if err != nil {
		// shouldn't happen because the tls.X509KeyPair call should catch this
		return err
	}

	ca.cert = cert
	ca.certPem = caCertPem
	ca.key = x509KeyPair.PrivateKey.(*rsa.PrivateKey)
	ca.keyPem = caKeyPem
	ca.tlsKeyPair = x509KeyPair
	ca.serial = ca.cert.SerialNumber

	return nil
}

func (ca *ca) writeCAToCache() error {
	if ca.cacheDir != "" {
		if err := os.WriteFile(filepath.Join(ca.cacheDir, caCertFile), ca.certPem, 0600); err != nil {
			return fmt.Errorf("cert write %w", err)
		}
		if err := os.WriteFile(filepath.Join(ca.cacheDir, caKeyFile), ca.keyPem, 0600); err != nil {
			return fmt.Errorf("key write %w", err)
		}
	}
	return nil
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

func (ca *ca) LoadOrCreateServerCert() (cert *tls.Certificate, fromCache bool, err error) {
	serverNum := ca.serverCount
	ca.serverCount++
	cachedCert, err := ca.loadServerCertFromCache(serverNum)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			// log and regenerate the server cert
			fmt.Printf("Error loading from cache: %v", err)
		}
	}
	if cachedCert != nil {
		return cachedCert, true, nil
	}

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
		return nil, false, fmt.Errorf("keygen: %w", err)
	}
	serverCertBytes, err := x509.CreateCertificate(rand.Reader, serverClaims, ca.cert, serverKey.Public(), ca.key)
	if err != nil {
		return nil, false, fmt.Errorf("sign: %w", err)
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
		return nil, false, fmt.Errorf("keypair: %w", err)
	}
	if err := ca.writeServerCertToCache(serverNum, serverKeyPem, serverCertPem); err != nil {
		fmt.Printf("Error writing to cache: %v", err)
	}
	return &serverKeyPair, false, nil
}

func (ca *ca) loadServerCertFromCache(num int) (*tls.Certificate, error) {
	if ca.cacheReadFailed || ca.cacheDir == "" {
		return nil, nil // don't read from the cache
	}
	// attempt to load the ca
	keyPem, err := os.ReadFile(filepath.Join(ca.cacheDir, fmt.Sprintf(serverKeyFile, num)))
	if err != nil {
		return nil, fmt.Errorf("server %v key %w", num, err)
	}
	certPem, err := os.ReadFile(filepath.Join(ca.cacheDir, fmt.Sprintf(serverCertFile, num)))
	if err != nil {
		return nil, fmt.Errorf("server %v cert %w", num, err)
	}
	x509KeyPair, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		return nil, fmt.Errorf("server %v keypair parse %w", num, err)
	}
	return &x509KeyPair, nil
}

func (ca *ca) writeServerCertToCache(num int, keyPem, certPem []byte) error {
	if ca.cacheDir == "" {
		return nil
	}
	if err := os.WriteFile(filepath.Join(ca.cacheDir, fmt.Sprintf(serverKeyFile, num)), keyPem, 0600); err != nil {
		return fmt.Errorf("write server %v key %w", num, err)
	}
	if err := os.WriteFile(filepath.Join(ca.cacheDir, fmt.Sprintf(serverCertFile, num)), certPem, 0600); err != nil {
		return fmt.Errorf("write server %v cert %w", num, err)
	}
	return nil
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

type caOption func(ca *ca)

func withCacheDir(cacheDir string) caOption {
	return func(ca *ca) {
		ca.cacheDir = cacheDir
	}
}
