// Package main is a simple program that connects to a local playground instance (on the default insecure port)
// and logs changes to demand and mode for a single electric device.
// It will also cycle through the device's modes every couple of minutes.
// The device name to use is specified by the first command line argument. Defaults to ELEC-001
package main

import (
	"context"
	"fmt"
	"github.com/smart-core-os/sc-api/go/traits"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:23557", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	deviceName := "ELEC-001"
	if len(os.Args) > 1 {
		deviceName = os.Args[1]
	}

	elec := traits.NewElectricApiClient(conn)

	group, ctx := errgroup.WithContext(context.Background())
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

func listenDemand(ctx context.Context, client traits.ElectricApiClient, name string) error {
	stream, err := client.PullDemand(ctx, &traits.PullDemandRequest{
		Name: name,
	})
	if err != nil {
		return err
	}
	defer stream.CloseSend()

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

	ticker := time.NewTicker(period)
	defer ticker.Stop()

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
