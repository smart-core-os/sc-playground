package mode

import (
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/router"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/mode"
	"github.com/smart-core-os/sc-playground/pkg/node"
)

func Activate(n *node.Node) {
	n.AddRouter(mode.NewApiRouter(
		mode.WithModeApiClientFactory(func(name string) (traits.ModeApiClient, error) {
			return mode.WrapApi(mode.NewModelServer(mode.NewModel())), nil
		}),
		router.WithOnChange(func(change router.Change) {
			n.Announce(change.Name, node.HasTrait(trait.Mode))
		}),
	))
}
