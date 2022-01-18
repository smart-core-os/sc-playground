package playground

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	minPort = 10000
	maxPort = 11000
)

var grpcPort = minPort + rand.Intn(maxPort-minPort)
var httpPort = minPort + rand.Intn(maxPort-minPort)
var httpsPort = minPort + rand.Intn(maxPort-minPort)

func TestRun(t *testing.T) {
	grpcPortStr := fmt.Sprintf(":%d", grpcPort)
	grpcBind = &grpcPortStr
	httpPortStr := fmt.Sprintf(":%d", httpPort)
	httpBind = &httpPortStr
	httpsPortStr := fmt.Sprintf(":%d", httpsPort)
	httpsBind = &httpsPortStr

	*mTlsEnabled = true

	ctx, done := context.WithCancel(context.Background())
	t.Cleanup(done)
	// run the server
	go func() {
		if err := runCtx(ctx); err != nil {
			t.Error(err)
		}
	}()

	// wait for the server to start (most of the time is spent generating the certs)
	if err := waitForServerStart(httpPort, 15*time.Second); err != nil {
		t.Fatal(fmt.Errorf("waiting for server start: %w", err))
	}

	// setup client tls
	// caCert := readCaCertFromFile(t, caCertFile)
	caCert := readCaCertFromUrl(t, fmt.Sprintf("http://localhost:%v/__/playground/ca-cert.pem", httpPort))
	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCert) {
		t.Fatal(fmt.Errorf("credentials: failed to append certificates"))
	}

	clientTlsConfig := &tls.Config{
		RootCAs: caPool,
	}
	if *mTlsEnabled {
		clientCert := readClientTlsFromUrl(t, fmt.Sprintf("http://localhost:%v/__/playground/client.pem", httpPort))
		clientTlsConfig.Certificates = []tls.Certificate{clientCert}
	}

	// setup a grpc client
	t.Run("grpc-secure", func(t *testing.T) {
		testGrpcCall(t, ctx, clientTlsConfig)
	})
	t.Run("grpc-insecure", func(t *testing.T) {
		t.SkipNow() // the server only supports secured connections
		testGrpcCall(t, ctx, nil)
	})

	// perform a http call
	t.Run("http/1", func(t *testing.T) {
		get(t, http.DefaultClient, fmt.Sprintf("http://localhost:%v/health", httpPort))
	})

	// perform a https 1.1 call
	t.Run("https/1", func(t *testing.T) {
		httpsClient := &http.Client{}
		httpsClient.Transport = &http.Transport{
			TLSClientConfig: clientTlsConfig,
		}
		get(t, httpsClient, fmt.Sprintf("https://localhost:%v/health", httpsPort))
	})

	// perform a https 2 call
	t.Run("https/2", func(t *testing.T) {
		https2Client := &http.Client{}
		https2Client.Transport = &http2.Transport{
			TLSClientConfig: clientTlsConfig,
		}
		get(t, https2Client, fmt.Sprintf("https://localhost:%v/health", httpsPort))
	})
}

// waitForServerStart waits for the given duration until localhost:{port} accepts a connection
func waitForServerStart(port int, d time.Duration) (err error) {
	now := time.Now()
	for deadline := now.Add(d); time.Now().Before(deadline); now = time.Now() {
		var conn net.Conn
		conn, err = net.DialTimeout("tcp", fmt.Sprintf("localhost:%v", port), deadline.Sub(now))
		if err == nil {
			_ = conn.Close()
			return nil
		}
	}
	return err
}

func readCaCertFromFile(t *testing.T, caCertFile string) []byte {
	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		t.Fatal(err)
	}
	return caCert
}

func readCaCertFromUrl(t *testing.T, caCertUrl string) []byte {
	resp, err := http.Get(caCertUrl)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}
func readClientTlsFromUrl(t *testing.T, url string) tls.Certificate {
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var block, keyBlock, certBlock *pem.Block
	for len(bytes) > 0 && (keyBlock == nil || certBlock == nil) {
		block, bytes = pem.Decode(bytes)
		if block == nil {
			break
		}
		switch block.Type {
		case "CERTIFICATE":
			certBlock = block
		case "RSA PRIVATE KEY":
			keyBlock = block
		}
	}
	if keyBlock == nil {
		t.Fatal("Missing client key block")
	}
	if certBlock == nil {
		t.Fatal("Missing client cert block")
	}

	keyPair, err := tls.X509KeyPair(pem.EncodeToMemory(certBlock), pem.EncodeToMemory(keyBlock))
	if err != nil {
		t.Fatal(err)
	}
	return keyPair
}

func testGrpcCall(t *testing.T, ctx context.Context, clientTlsConfig *tls.Config) {
	var options []grpc.DialOption
	if clientTlsConfig == nil {
		options = append(options, grpc.WithInsecure())
	} else {
		clientTls := credentials.NewTLS(clientTlsConfig)
		options = append(options, grpc.WithTransportCredentials(clientTls))
	}
	conn, err := grpc.Dial(
		fmt.Sprintf("localhost:%d", grpcPort),
		options...)
	if err != nil {
		t.Fatal(err)
	}

	powerSupplyClient := traits.NewPowerSupplyApiClient(conn)
	capacity, err := powerSupplyClient.GetPowerCapacity(ctx, &traits.GetPowerCapacityRequest{
		Name: "TST-001",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("capacity %v", capacity)
}

func get(t *testing.T, client *http.Client, url string) {
	got, err := client.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer got.Body.Close()
	body, err := io.ReadAll(got.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Got %v %v %v: %s", url, got.Proto, got.Status, body)
}
