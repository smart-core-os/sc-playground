package apis

import (
	"log"
	"math/rand"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/server"
	"github.com/smart-core-os/sc-golang/pkg/trait/energystorage"
	sim "github.com/smart-core-os/sc-playground/internal/simulated/energystorage"
)

func EnergyStorageApi() server.GrpcApi {
	r := energystorage.NewApiRouter(
		energystorage.WithEnergyStorageApiClientFactory(func(name string) (traits.EnergyStorageApiClient, error) {
			log.Printf("Creating EnergyStorageClient(%v)", name)
			model := energystorage.NewModel()

			randStart := time.Duration(rand.Int63n(int64(10 * time.Minute)))
			oscillator := sim.NewOscillator(model, sim.WithCycleStart(time.Now().Add(-randStart)))
			go func() {
				// loop forever
				ticker := time.NewTicker(250 * time.Millisecond)
				defer ticker.Stop()
				for {
					if err := oscillator.Scrub(<-ticker.C); err != nil {
						log.Printf("Error returned from [%v]oscillator.Scrub, %v", name, err)
					}
				}
			}()

			return energystorage.WrapApi(energystorage.NewModelServer(model, energystorage.ReadOnly())), nil
		}),
	)
	return r
}
