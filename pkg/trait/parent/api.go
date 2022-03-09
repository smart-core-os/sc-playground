package parent

import (
	"log"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/router"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/parent"
	"github.com/smart-core-os/sc-playground/pkg/node"
)

func Activate(n *node.Node) {
	n.AddRouter(parent.NewApiRouter(
		parent.WithParentApiClientFactory(func(name string) (traits.ParentApiClient, error) {
			return parent.WrapApi(parent.NewModelServer(parent.NewModel())), nil
		}),
		router.WithOnCommit(func(name string, client interface{}) {
			log.Printf("ParentApiClient(%v) auto-created", name)
			n.Announce(name, node.HasTrait(trait.Parent))
		}),
	))
}
