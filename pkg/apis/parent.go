package apis

import (
	"log"

	"github.com/smart-core-os/sc-golang/pkg/router"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/parent"
	parent2 "github.com/smart-core-os/sc-playground/pkg/apis/parent"
	"github.com/smart-core-os/sc-playground/pkg/apis/registry"
)

func ParentApi(traiter parent2.Traiter, adder registry.Adder) *parent.ApiRouter {
	r := parent.NewApiRouter(router.WithFactory(func(name string) (interface{}, error) {
		log.Printf("Creating ParentApiClient(%v)", name)
		traiter.Trait(name, trait.Parent)
		model := parent.NewModel()
		return parent.WrapApi(parent.NewModelServer(model)), nil
	}))

	adder.Add(registry.ParentApiRegistry{ApiRouter: r, Traiter: traiter})
	return r
}
