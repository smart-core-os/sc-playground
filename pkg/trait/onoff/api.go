package onoff

import (
	"log"
	"math/rand"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/router"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/onoff"
	"github.com/smart-core-os/sc-golang/pkg/wrap"
	"github.com/smart-core-os/sc-playground/pkg/node"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func Activate(n *node.Node) {
	r := onoff.NewApiRouter(
		onoff.WithOnOffApiClientFactory(func(name string) (traits.OnOffApiClient, error) {
			var onOrOff traits.OnOff_State
			n := rand.Intn(10)
			if n < 5 {
				onOrOff = traits.OnOff_OFF
			} else {
				onOrOff = traits.OnOff_ON
			}
			return onoff.WrapApi(onoff.NewModelServer(onoff.NewModel(onOrOff))), nil
		}),
		router.WithOnChange(func(change router.Change) {
			if !change.Auto {
				return
			}
			name := change.Name
			model, ok := wrap.UnwrapFully(change.New).(*onoff.Model)
			if !ok {
				return
			}

			currentVal, err := model.GetOnOff()
			if err != nil {
				log.Printf("OnOffApiClient(%v) auto-created %v", name, currentVal.State)
			} else {
				log.Printf("OnOffApiClient(%v) auto-created (%v)", name, err)
			}
			n.Announce(name, node.HasTrait(trait.OnOff))
		}),
	)
	n.AddRouter(r)
	n.AddTraitFactory(trait.OnOff, func(name string, _ proto.Message) error {
		_, err := r.Get(name)
		return err
	})
	n.AddClientFactory(trait.OnOff, func(conn *grpc.ClientConn) interface{} {
		return traits.NewOnOffApiClient(conn)
	})
}
