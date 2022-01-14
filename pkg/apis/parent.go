package apis

import (
	"log"

	"github.com/smart-core-os/sc-golang/pkg/router"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/parent"
)

func ParentApi(traiter Traiter) *parent.ApiRouter {
	return parent.NewApiRouter(router.WithFactory(func(name string) (interface{}, error) {
		log.Printf("Creating ParentApiClient(%v)", name)
		traiter.Trait(name, trait.Parent)
		model := parent.NewModel()
		return parent.WrapApi(parent.NewModelServer(model)), nil
	}))
}
