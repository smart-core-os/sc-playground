package airtemperature

import (
	"context"
	"log"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/router"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/airtemperature"
	"github.com/smart-core-os/sc-golang/pkg/wrap"
	"github.com/smart-core-os/sc-playground/pkg/node"
)

func Activate(n *node.Node) {
	n.AddRouter(airtemperature.NewApiRouter(
		airtemperature.WithAirTemperatureApiClientFactory(func(name string) (traits.AirTemperatureApiClient, error) {
			return airtemperature.WrapApi(airtemperature.NewMemoryDevice()), nil
		}),
		router.WithOnCommit(func(name string, client interface{}) {
			model, ok := wrap.UnwrapFully(client).(*airtemperature.MemoryDevice)
			if !ok {
				return
			}

			currentVal, err := model.GetAirTemperature(context.Background(), &traits.GetAirTemperatureRequest{})
			if err != nil {
				log.Printf("AirTemperatureApiClient(%v) auto-created (%v)", name, err)
			} else {
				log.Printf("AirTemperatureApiClient(%v) auto-created %v", name, currentVal.TemperatureGoal)
			}
			n.Announce(name, node.HasTrait(trait.AirTemperature))
		}),
	))
}
