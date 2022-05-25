package playground

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/smart-core-os/sc-playground/pkg/config"
	"github.com/smart-core-os/sc-playground/pkg/device/evcharger"
	"github.com/smart-core-os/sc-playground/pkg/node"
	"github.com/smart-core-os/sc-playground/pkg/playpb"
	"github.com/smart-core-os/sc-playground/pkg/run"
	"github.com/smart-core-os/sc-playground/pkg/sim"
	"github.com/smart-core-os/sc-playground/pkg/sim/boot"
	"github.com/smart-core-os/sc-playground/pkg/sim/stats"
	"github.com/smart-core-os/sc-playground/pkg/trait/airtemperature"
	"github.com/smart-core-os/sc-playground/pkg/trait/booking"
	"github.com/smart-core-os/sc-playground/pkg/trait/electric"
	"github.com/smart-core-os/sc-playground/pkg/trait/energystorage"
	"github.com/smart-core-os/sc-playground/pkg/trait/fanspeed"
	"github.com/smart-core-os/sc-playground/pkg/trait/light"
	"github.com/smart-core-os/sc-playground/pkg/trait/metadata"
	"github.com/smart-core-os/sc-playground/pkg/trait/mode"
	"github.com/smart-core-os/sc-playground/pkg/trait/occupancysensor"
	"github.com/smart-core-os/sc-playground/pkg/trait/onoff"
	"github.com/smart-core-os/sc-playground/pkg/trait/openclose"
	"github.com/smart-core-os/sc-playground/pkg/trait/parent"
	"github.com/smart-core-os/sc-playground/pkg/trait/powersupply"
	"github.com/smart-core-os/sc-playground/ui"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"
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

	deviceConfig = flag.String("devices", "", "Configuration describing which devices to load into the playground on boot")
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

	serverDeviceName := "scos/play"
	rootNode := newRootNode(serverDeviceName)

	// setup the playground api
	rootApi := playpb.New(rootNode)

	// Start the simulation loop
	simulation, err := boot.CreateSimulation(rootNode, sim.WithFramer(stats.CFramer(60*5, rootApi)))
	if err != nil {
		return err
	}
	go func() {
		err := simulation.Run(ctx)
		if err != nil {
			log.Printf("Simulation ended: %v", err)
		}
	}()

	// setup configured devices
	if err := seedDevices(*deviceConfig, rootNode); err != nil {
		return fmt.Errorf("--devices %v : %w", *deviceConfig, err)
	}

	return run.Serve(
		run.WithContext(ctx),
		run.WithDefaultName(serverDeviceName),
		run.WithApis(rootNode, rootApi),
		run.WithGrpcAddress(*grpcBind),
		run.WithHttpAddress(*httpBind),
		run.WithHttpsAddress(*httpsBind),
		run.WithHostedFS(ui.Playground),
		run.WithHostedFSNotFound(ui.NotFound),
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
	// trait apis
	parent.Activate(root)
	metadata.Activate(root)

	airtemperature.Activate(root)
	booking.Activate(root)
	electric.Activate(root)
	energystorage.Activate(root)
	fanspeed.Activate(root)
	light.Activate(root)
	mode.Activate(root)
	occupancysensor.Activate(root)
	onoff.Activate(root)
	openclose.Activate(root)
	powersupply.Activate(root)

	// device apis
	evcharger.Activate(root)
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

func seedDevices(deviceConfig string, rootNode *node.Node) error {
	devices := new(config.File)
	if deviceConfig != "" {
		if err := config.ReadFile(deviceConfig, devices); err != nil {
			return err
		}
	}
	var err error
	log.Printf("Seeding %v devices from config", len(devices.Devices))
	for name, device := range devices.Devices {
		if device.IsLocal() {
			for _, trait := range device.Traits {
				traitConfig, err2 := trait.Config.UnmarshalProto(proto.UnmarshalOptions{})
				if err2 != nil {
					err = multierr.Append(err, err2)
					continue
				}
				err2 = rootNode.CreateDeviceTrait(name, trait.Name, traitConfig)
				if err2 != nil {
					err = multierr.Append(err, err2)
					continue
				}
			}
		} else {
			// we should probably cache the connections here, but just want to get it running for now
			var endpoint config.Node
			if device.Node.NameOnly {
				var ok bool
				endpoint, ok = devices.Nodes[device.Node.Name]
				if !ok {
					return fmt.Errorf("devices[%v].node[%v] not found", name, device.Node.Name)
				}
			} else {
				endpoint = device.Node.Node
			}

			conn, err := endpoint.ResolveRemoteConn(context.Background(), rootNode)
			if err != nil {
				return fmt.Errorf("devices[%v] dial %w", name, err)
			}
			var err2 error
			var features []node.Feature
			for _, traitName := range device.Traits {
				client, err := rootNode.CreateTraitClient(traitName.Name, conn)
				if err != nil {
					err2 = multierr.Append(err2, err)
					continue
				}
				features = append(features, node.HasTrait(traitName.Name, node.WithClients(client), node.NoAddMetadata()))
			}
			rootNode.Announce(name, features...)
			if err2 != nil {
				err = multierr.Append(err, err2)
			}
		}
	}
	return err
}
