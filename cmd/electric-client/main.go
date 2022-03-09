// Package main is a simple program that connects to a local playground instance (on the default insecure port)
// and logs changes to demand and mode for a single electric device.
// It will also cycle through the device's modes every couple of minutes.
// The device name to use is specified by the first command line argument. Defaults to ELEC-001
package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
)

var (
	node     = flag.String("node", "localhost:23557", "The node host:port to connect to")
	useTLS   = flag.Bool("tls", false, "When specified, use the local machines certificate store for server trust")
	useTLSCA = flag.String("tls-ca", "", "A path to a CA certification in the chain of trust for the server")
)

func main() {
	flag.Parse()

	log.Printf("Connecting to %v\n", *node)
	ctx, done := context.WithCancel(context.Background())
	defer done()

	dialCtx, dialDone := context.WithTimeout(ctx, 5*time.Second)
	defer dialDone()
	var opts []grpc.DialOption
	if *useTLS {
		c := &tls.Config{}
		if *useTLSCA != "" {
			var err error
			c.RootCAs, err = readCACert()
			if err != nil {
				log.Fatalf("--tls-ca=%v %v", *useTLSCA, err)
			}
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(c)))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.DialContext(dialCtx, *node, opts...)
	if err != nil {
		panic(err)
	}

	go func() {
		var state connectivity.State
		for conn.WaitForStateChange(ctx, state) {
			state = conn.GetState()
			log.Printf("Client connection %v", state)
		}
	}()

	deviceName := "ELEC-001"
	if len(flag.Args()) > 0 {
		deviceName = flag.Arg(0)
	}

	elec := traits.NewElectricApiClient(conn)

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		err := listenDemand(ctx, elec, deviceName)
		if err != nil {
			log.Println("listenDemand finished with error", err)
		}
		return err
	})
	group.Go(func() error {
		err := listenActiveMode(ctx, elec, deviceName)
		if err != nil {
			log.Println("listenActiveMode finished with error", err)
		}
		return err
	})
	group.Go(func() error {
		err := changeModes(ctx, elec, deviceName)
		if err != nil {
			log.Println("changeModes finished with error", err)
		}
		return err
	})

	err = group.Wait()
	if err != nil {
		panic(err)
	}
}

func readCACert() (*x509.CertPool, error) {
	path := *useTLSCA
	var caPem []byte
	if strings.HasPrefix(path, "http") {
		response, err := http.Get(path)
		if err != nil {
			return nil, fmt.Errorf("request %w", err)
		}
		if response.StatusCode < 200 || response.StatusCode >= 400 {
			return nil, fmt.Errorf("status %v %v", response.StatusCode, response.Status)
		}
		defer response.Body.Close()
		pem, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("read body %w", err)
		}
		caPem = pem
	} else {
		var err error
		caPem, err = os.ReadFile(*useTLSCA)
		if err != nil {
			return nil, fmt.Errorf("ca-certfile %w", err)
		}
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caPem) {
		return nil, fmt.Errorf("ca-certfile failed to add to cert pool")
	}
	return pool, nil
}

func listenDemand(ctx context.Context, client traits.ElectricApiClient, name string) error {
	stream, err := client.PullDemand(ctx, &traits.PullDemandRequest{
		Name: name,
	})
	if err != nil {
		return err
	}
	defer stream.CloseSend()

	initial, err := client.GetDemand(ctx, &traits.GetDemandRequest{
		Name: name,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Initial demand is %f\n", initial.Current)

	for {
		res, err := stream.Recv()
		if err != nil {
			return err
		}

		for _, change := range res.Changes {
			fmt.Printf("At %v: Demand is now %f\n", change.ChangeTime.AsTime(), change.Demand.Current)
		}
	}
}

func listenActiveMode(ctx context.Context, client traits.ElectricApiClient, name string) error {
	stream, err := client.PullActiveMode(ctx, &traits.PullActiveModeRequest{
		Name: name,
	})
	if err != nil {
		return err
	}
	defer stream.CloseSend()

	initial, err := client.GetActiveMode(ctx, &traits.GetActiveModeRequest{
		Name: name,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Initial active mode is %q - %q\n", initial.Id, initial.Title)

	for {
		res, err := stream.Recv()
		if err != nil {
			return err
		}

		for _, change := range res.Changes {
			fmt.Printf("At %v: Active mode is now %s\n", change.ChangeTime.AsTime(), change.ActiveMode.Id)
		}
	}
}

func changeModes(ctx context.Context, client traits.ElectricApiClient, name string) error {
	const period = 2 * time.Minute

	modes, err := listModes(ctx, client, name)
	if err != nil {
		return err
	}

	if len(modes) == 0 {
		// don't change any modes
		log.Printf("Device %q has no modes", name)
		<-ctx.Done()
		return ctx.Err()
	}

	ticker := time.NewTicker(period)
	defer ticker.Stop()
	force := make(chan struct{})

	go func() {
		in := bufio.NewReader(os.Stdin)
		fmt.Println("Press enter to switch modes...")
		for {
			_, _ = in.ReadString('\n')

			select {
			case <-ctx.Done():
				return
			case force <- struct{}{}:
			}
		}
	}()

	i := 0
	for {
		mode := modes[i]

		log.Printf("Activating mode %q with id %s", mode.Title, mode.Id)

		_, err = client.UpdateActiveMode(ctx, &traits.UpdateActiveModeRequest{
			Name:       name,
			ActiveMode: &traits.ElectricMode{Id: mode.Id},
		})
		if err != nil {
			return err
		}

		// prepare for next loop
		select {
		case <-ticker.C:
		case <-force:
			ticker.Stop() // if you're doing it by hand, disable auto switching
		case <-ctx.Done():
			return ctx.Err()
		}
		i++
		if i >= len(modes) {
			i = 0
		}
	}

}

func listModes(ctx context.Context, client traits.ElectricApiClient, deviceName string) ([]*traits.ElectricMode, error) {
	var modes []*traits.ElectricMode

	token := ""
	for {
		res, err := client.ListModes(ctx, &traits.ListModesRequest{
			Name:      deviceName,
			PageToken: token,
		})
		if err != nil {
			return nil, err
		}

		modes = append(modes, res.Modes...)

		token = res.NextPageToken
		if token == "" {
			break
		}
	}

	return modes, nil
}
