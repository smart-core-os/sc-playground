package energystorage

import (
	"log"
	"math/rand"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/energystorage"
	sim "github.com/smart-core-os/sc-playground/internal/simulated/energystorage"
	"github.com/smart-core-os/sc-playground/pkg/node"
)

func Activate(n *node.Node) {
	r := energystorage.NewApiRouter(
		energystorage.WithEnergyStorageApiClientFactory(func(name string) (traits.EnergyStorageApiClient, error) {
			log.Printf("Creating EnergyStorageClient(%v)", name)
			model := energystorage.NewModel()

			randStart := time.Duration(rand.Int63n(int64(10 * time.Minute)))
			oscillator := sim.NewOscillator(model, sim.WithCycleStart(time.Now().Add(-randStart)))

			t := time.Now()
			t = t.Round(15 * time.Second)

			n.Announce(name, node.HasSimulation(oscillator), node.HasSimulation(oscillator), node.HasTrait(trait.EnergyStorage))
			return energystorage.WrapApi(energystorage.NewModelServer(model, energystorage.ReadOnly())), nil
		}),
	)
	n.AddRouter(r)
}
