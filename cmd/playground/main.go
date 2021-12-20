package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/smart-core-os/sc-playground/pkg/apis"
	"github.com/smart-core-os/sc-playground/pkg/run"
	"github.com/smart-core-os/sc-playground/ui"
)

var (
	grpcBind       = flag.String("bind", ":23557", "grpc server bind")
	httpBind       = flag.String("http", ":8080", "http admin and grpc-web binding")
	httpsBind      = flag.String("https", ":8443", "https admin and grpc-web binding")
	serverCertFile = flag.String("server-certfile", "", "a path to the servers cert file")
	serverKeyFile  = flag.String("server-keyfile", "", "a path to the servers private key file")
	insecure       = flag.Bool("insecure", false, "Do not generate ssl certificates when server certs are not specified")
	caCertFile     = flag.String("ca-certfile", "", "a path to the CA cert file. Used during mutual-tls to verify client connections")
	mTlsEnabled    = flag.Bool("mtls", false, "Enable mutual TLS. Use --ca-certfile to specify a CA cert that will have signed client certs."+
		"Or GET /__/playground/client.pem to download a new PEM encoded client cert and key pair based on the internal self-signed CA.")
)

func main() {
	if err := Run(); err != nil {
		log.Printf("Exiting: %v", err)
	}
}

func Run() error {
	flag.Parse()
	return runCtx(context.Background())
}

func runCtx(ctx context.Context) error {
	serverTlsConfig, err := serverTlsConfig()
	if err != nil {
		return err
	}

	return run.Serve(
		run.WithContext(ctx),
		run.WithApis(
			apis.BookingApi(),
			apis.ElectricApi(),
			apis.EnergyStorageApi(),
			apis.OccupancyApi(),
			apis.OnOffApi(),
			apis.PowerSupplyApi(),
		),
		run.WithGrpcAddress(*grpcBind),
		run.WithHttpAddress(*httpBind),
		run.WithHttpsAddress(*httpsBind),
		run.WithHostedFS(ui.Playground),
		withInsecure(),
		withMTLS(),
		run.WithGrpcTls(serverTlsConfig),
		run.WithHttpHealth("/health"),
	)
}

func serverTlsConfig() (*tls.Config, error) {
	if *serverCertFile != "" && *serverKeyFile != "" {
		keyPair, err := tls.LoadX509KeyPair(*serverCertFile, *serverKeyFile)
		if err != nil {
			return nil, err
		}
		config := &tls.Config{
			Certificates: []tls.Certificate{keyPair},
		}
		if *mTlsEnabled {
			if *caCertFile == "" {
				return nil, fmt.Errorf("with mtls enabled, missing ca-certfile")
			}
			caPem, err := os.ReadFile(*caCertFile)
			if err != nil {
				return nil, fmt.Errorf("ca-certfile: %w", err)
			}
			pool := x509.NewCertPool()
			if !pool.AppendCertsFromPEM(caPem) {
				return nil, fmt.Errorf("ca-certfile failed to add to cert pool")
			}
			config.RootCAs = pool
			config.ClientAuth = tls.RequireAndVerifyClientCert
		}
		return config, nil
	} else if *serverCertFile != "" || *serverKeyFile != "" {
		return nil, fmt.Errorf("both server-certfile and server-keyfile must be preset or absent")
	}
	return nil, nil
}

func withInsecure() run.ConfigOption {
	if *insecure {
		return run.WithInsecure()
	} else {
		return run.NilConfigOption
	}
}

func withMTLS() run.ConfigOption {
	if *mTlsEnabled {
		return run.WithMTLS()
	} else {
		return run.NilConfigOption
	}
}
