package energystorage

import (
	"log"
	"math/rand"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/router"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/energystorage"
	"github.com/smart-core-os/sc-golang/pkg/wrap"
	sim "github.com/smart-core-os/sc-playground/internal/simulated/energystorage"
	"github.com/smart-core-os/sc-playground/pkg/node"
	"google.golang.org/protobuf/proto"
)

func Activate(n *node.Node) {
	r := energystorage.NewApiRouter(
		energystorage.WithEnergyStorageApiClientFactory(func(name string) (traits.EnergyStorageApiClient, error) {
			return energystorage.WrapApi(energystorage.NewModelServer(energystorage.NewModel(), energystorage.ReadOnly())), nil
		}),
		router.WithOnCommit(func(name string, client interface{}) {
			log.Printf("EnergyStorageApiClient(%v) auto-created", name)
			model, ok := wrap.UnwrapFully(client).(*energystorage.Model)
			if !ok {
				return
			}

			randStart := time.Duration(rand.Int63n(int64(10 * time.Minute)))
			oscillator := sim.NewOscillator(model, sim.WithCycleStart(time.Now().Add(-randStart)))

			t := time.Now()
			t = t.Round(15 * time.Second)

			n.Announce(name, node.HasSimulation(oscillator), node.HasSimulation(oscillator), node.HasTrait(trait.EnergyStorage))
		}),
	)
	n.AddRouter(r)
	n.AddTraitFactory(trait.EnergyStorage, func(name string, _ proto.Message) error {
		_, err := r.Get(name)
		return err
	})
}
