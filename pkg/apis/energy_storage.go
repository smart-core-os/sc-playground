package apis

import (
	"log"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/server"
	"github.com/smart-core-os/sc-golang/pkg/trait/energystorage"
)

func EnergyStorageApi() server.GrpcApi {
	r := energystorage.NewRouter()
	r.Factory = func(name string) (traits.EnergyStorageApiClient, error) {
		log.Printf("Creating EnergyStorageClient(%v)", name)
		model := energystorage.NewModel()
		return energystorage.Wrap(energystorage.NewModelServer(model)), nil
	}
	return r
}
