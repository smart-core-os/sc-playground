package apis

import (
	"log"
	"math/rand"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/server"
	"github.com/smart-core-os/sc-golang/pkg/trait/onoff"
)

func OnOffApi() server.GrpcApi {
	r := onoff.NewRouter()
	r.Factory = func(name string) (traits.OnOffApiClient, error) {
		var onOrOff traits.OnOff_State
		n := rand.Intn(10)
		if n < 5 {
			onOrOff = traits.OnOff_OFF
		} else {
			onOrOff = traits.OnOff_ON
		}
		log.Printf("Creating OnOffClient(%v)=%v", name, onOrOff)
		return onoff.Wrap(onoff.NewMemoryDevice(onOrOff)), nil
	}
	return r
}
