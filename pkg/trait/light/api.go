package light

import (
	"context"
	"log"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/router"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/light"
	"github.com/smart-core-os/sc-golang/pkg/wrap"
	"github.com/smart-core-os/sc-playground/pkg/node"
)

func Activate(n *node.Node) {
	n.AddRouter(light.NewApiRouter(
		light.WithLightApiClientFactory(func(name string) (traits.LightApiClient, error) {
			return light.WrapApi(light.NewMemoryDevice()), nil
		}),
		router.WithOnCommit(func(name string, client interface{}) {
			model, ok := wrap.UnwrapFully(client).(*light.MemoryDevice)
			if !ok {
				return
			}

			currentVal, err := model.GetBrightness(context.Background(), &traits.GetBrightnessRequest{})
			if err != nil {
				log.Printf("LightApiClient(%v) auto-created (%v)", name, err)

			} else {
				log.Printf("LightApiClient(%v) auto-created %v", name, currentVal.LevelPercent)
			}
			n.Announce(name, node.HasTrait(trait.Light))
		}),
	))
}
