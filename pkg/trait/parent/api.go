package parent

import (
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/parent"
	"github.com/smart-core-os/sc-playground/pkg/node"
)

func Activate(n *node.Node) {
	n.AddRouter(parent.NewApiRouter(parent.WithParentApiClientFactory(func(name string) (traits.ParentApiClient, error) {
		n.Announce(name, node.HasTrait(trait.Parent))
		return parent.WrapApi(parent.NewModelServer(parent.NewModel())), nil
	})))
}
