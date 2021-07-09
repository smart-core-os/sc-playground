package apis

import (
	"log"
	"math/rand"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/memory"
	"github.com/smart-core-os/sc-golang/pkg/router"
	"github.com/smart-core-os/sc-golang/pkg/server"
	"github.com/smart-core-os/sc-golang/pkg/wrap"
)

func OnOffApi() server.GrpcApi {
	r := router.NewOnOffApiRouter()
	r.Factory = func(name string) (traits.OnOffApiClient, error) {
		var onOrOff traits.OnOff_State
		n := rand.Intn(10)
		if n < 5 {
			onOrOff = traits.OnOff_OFF
		} else {
			onOrOff = traits.OnOff_ON
		}
		log.Printf("Creating OnOffClient(%v)=%v", name, onOrOff)
		return wrap.OnOffApiServer(memory.NewOnOffApi(onOrOff)), nil
	}
	return r
}
