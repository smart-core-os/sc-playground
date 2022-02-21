package onoff

import (
	"log"
	"math/rand"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/onoff"
	"github.com/smart-core-os/sc-playground/pkg/node"
)

func Activate(n *node.Node) {
	n.AddRouter(onoff.NewApiRouter(onoff.WithOnOffApiClientFactory(func(name string) (traits.OnOffApiClient, error) {
		n.Announce(name, node.HasTrait(trait.OnOff))
		var onOrOff traits.OnOff_State
		n := rand.Intn(10)
		if n < 5 {
			onOrOff = traits.OnOff_OFF
		} else {
			onOrOff = traits.OnOff_ON
		}
		log.Printf("Creating OnOffClient(%v)=%v", name, onOrOff)
		return onoff.WrapApi(onoff.NewModelServer(onoff.NewModel(onOrOff))), nil
	})))
}
