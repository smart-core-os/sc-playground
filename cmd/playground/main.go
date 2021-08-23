package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"

	"github.com/smart-core-os/sc-playground/pkg/apis"
	"github.com/smart-core-os/sc-playground/pkg/run"
	"github.com/smart-core-os/sc-playground/ui"
)

var (
	grpcBind    = flag.String("bind", ":23557", "grpc server bind")
	httpBind    = flag.String("http", ":8080", "http admin and grpc-web binding")
	httpsBind   = flag.String("https", ":8443", "https admin and grpc-web binding")
	serveFolder = flag.String("public", "", "a folder to host")
	// caCertFile = flag.String("ca-certfile", "", "a path to the CA cert file")
	serverCertFile = flag.String("server-certfile", "", "a path to the servers cert file")
	serverKeyFile  = flag.String("server-keyfile", "", "a path to the servers private key file")
	insecure       = flag.Bool("insecure", false, "Do not generate ssl certificates when server certs are not specified")
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
	tlsConfig, err := tlsConfig()
	if err != nil {
		return err
	}

	return run.Serve(
		run.WithContext(ctx),
		run.WithApis(
			apis.BookingApi(),
			apis.OccupancyApi(),
			apis.OnOffApi(),
			apis.PowerSupplyApi(),
		),
		run.WithGrpcAddress(*grpcBind),
		run.WithHttpAddress(*httpBind),
		run.WithHttpsAddress(*httpsBind),
		withHostedOrEmbedded(),
		run.WithGrpcTls(tlsConfig),
		run.WithHttpHealth("/health"),
	)
}

func tlsConfig() (*tls.Config, error) {
	if *serverCertFile != "" && *serverKeyFile != "" {
		keyPair, err := tls.LoadX509KeyPair(*serverCertFile, *serverKeyFile)
		if err != nil {
			return nil, err
		}
		return &tls.Config{
			Certificates: []tls.Certificate{keyPair},
		}, nil
	} else if *serverCertFile != "" || *serverKeyFile != "" {
		return nil, fmt.Errorf("both server-certfile and server-keyfile must be preset or absent")
	}
	return nil, nil
}

func withHostedOrEmbedded() run.ConfigOption {
	if *serveFolder != "" {
		return run.WithHostedDir(*serveFolder)
	}
	return run.WithHostedFS(ui.Playground)
}
