package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"testing"

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

	*serverCertFile = "/Users/matt/Projects/smartcore/sc-playground/scripts/server.crt"
	*serverKeyFile = "/Users/matt/Projects/smartcore/sc-playground/scripts/server.key"
	caCertFile := "/Users/matt/Projects/smartcore/sc-playground/scripts/ca.crt"

	ctx, done := context.WithCancel(context.Background())
	t.Cleanup(done)
	// run the server
	go func() {
		if err := runCtx(ctx); err != nil {
			t.Error(err)
		}
	}()

	// setup client tls
	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		t.Fatal(err)
	}
	clientCerts := x509.NewCertPool()
	if !clientCerts.AppendCertsFromPEM(caCert) {
		t.Fatal(fmt.Errorf("credentials: failed to append certificates"))
	}
	clientTlsConfig := &tls.Config{
		RootCAs: clientCerts,
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
