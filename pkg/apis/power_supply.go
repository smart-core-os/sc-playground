package apis

import (
	"math"
	"math/rand"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/memory"
	"github.com/smart-core-os/sc-golang/pkg/router"
	"github.com/smart-core-os/sc-golang/pkg/server"
	"github.com/smart-core-os/sc-golang/pkg/wrap"
)

func PowerSupplyApi() server.GrpcApi {
	r := router.NewPowerSupplyApiRouter()
	// todo: adjust load over time
	r.Factory = func(name string) (traits.PowerSupplyApiClient, error) {
		device := memory.NewPowerSupplyApi()
		// seed with a random load
		device.SetLoad(float32(math.Round(rand.Float64()*40*100) / 100))
		return wrap.PowerSupplyApiServer(device), nil
	}
	return r
}
