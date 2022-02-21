package playground

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/smart-core-os/sc-playground/pkg/node"
	"github.com/smart-core-os/sc-playground/pkg/run"
	"github.com/smart-core-os/sc-playground/pkg/sim/boot"
	"github.com/smart-core-os/sc-playground/pkg/trait/booking"
	"github.com/smart-core-os/sc-playground/pkg/trait/electric"
	"github.com/smart-core-os/sc-playground/pkg/trait/energystorage"
	"github.com/smart-core-os/sc-playground/pkg/trait/metadata"
	"github.com/smart-core-os/sc-playground/pkg/trait/occupancysensor"
	"github.com/smart-core-os/sc-playground/pkg/trait/onoff"
	"github.com/smart-core-os/sc-playground/pkg/trait/parent"
	"github.com/smart-core-os/sc-playground/pkg/trait/powersupply"
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
	forceCertGen = flag.Bool("force-cert-gen", false, "Force the generation of certificates, do not use cached certs. Ignored when specifying cert files")
	certCacheDir = flag.String("cert-cache-dir", "", "Generated certificates will be placed here and loaded from here unless --force-cert-gen or cert files are specified. "+
		"Defaults to a directory in TMP.")
)

func Main() {
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

	serverDeviceName := "scos/apps/playground"
	root := newRootNode(serverDeviceName)

	// Setup some devices to start us off
	simulation, err := boot.CreateSimulation(root)
	if err != nil {
		return err
	}
	go func() {
		err := simulation.Run(ctx)
		if err != nil {
			log.Printf("Simulation ended: %v", err)
		}
	}()

	return run.Serve(
		run.WithContext(ctx),
		run.WithDefaultName(serverDeviceName),
		run.WithApis(root),
		run.WithGrpcAddress(*grpcBind),
		run.WithHttpAddress(*httpBind),
		run.WithHttpsAddress(*httpsBind),
		run.WithHostedFS(ui.Playground),
		withInsecure(),
		withMTLS(),
		run.WithGrpcTls(serverTlsConfig),
		withForceCertGen(),
		run.WithCertCacheDir(*certCacheDir),
		run.WithHttpHealth("/health"),
	)
}

func newRootNode(serverDeviceName string) *node.Node {
	root := node.New(serverDeviceName)
	booking.Activate(root)
	electric.Activate(root)
	energystorage.Activate(root)
	metadata.Activate(root)
	occupancysensor.Activate(root)
	onoff.Activate(root)
	parent.Activate(root)
	powersupply.Activate(root)
	return root
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

func withForceCertGen() run.ConfigOption {
	if *forceCertGen {
		return run.WithForceCertGen()
	} else {
		return run.NilConfigOption
	}
}
