package apis

import (
	"log"
	"math/rand"

	"git.vanti.co.uk/smartcore/sc-api/go/traits"
	"git.vanti.co.uk/smartcore/sc-golang/pkg/memory"
	"git.vanti.co.uk/smartcore/sc-golang/pkg/router"
	"git.vanti.co.uk/smartcore/sc-golang/pkg/server"
	"git.vanti.co.uk/smartcore/sc-golang/pkg/wrap"
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
