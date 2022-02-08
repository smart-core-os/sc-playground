package apis

import (
	"log"
	"math/rand"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/server"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"
	"github.com/smart-core-os/sc-playground/pkg/apis/parent"
	"github.com/smart-core-os/sc-playground/pkg/apis/registry"
)

func OccupancyApi(traiter parent.Traiter, adder registry.Adder) server.GrpcApi {
	// handle random changes in occupancy
	var devices []struct {
		api  *occupancysensor.Model
		name string
	}
	go func() {
		for range time.Tick(30 * time.Second) {
			if len(devices) > 0 {
				i := rand.Intn(len(devices))
				o := randomOccupancy()
				log.Printf("SetOccupancy(%v) %v (people=%v)", devices[i].name, o.State, o.PeopleCount)
				_, _ = devices[i].api.SetOccupancy(o)
			}
		}
	}()

	r := occupancysensor.NewApiRouter(
		occupancysensor.WithOccupancySensorApiClientFactory(func(name string) (traits.OccupancySensorApiClient, error) {
			traiter.Trait(name, trait.OccupancySensor)
			initial := randomOccupancy()
			log.Printf("Creating OccupancyApiClient(%v) %v (people=%v)", name, initial.State, initial.PeopleCount)
			api := occupancysensor.NewModel(initial)
			devices = append(devices, struct {
				api  *occupancysensor.Model
				name string
			}{api: api, name: name})
			return occupancysensor.WrapApi(occupancysensor.NewModelServer(api)), nil
		}),
	)

	adder.Add(registry.OccupancySensorApiRegistry{ApiRouter: r, Traiter: traiter})
	return r
}

func randomOccupancy() *traits.Occupancy {
	var occupiedOrNot traits.Occupancy_State
	n := rand.Intn(10)
	if n < 5 {
		occupiedOrNot = traits.Occupancy_OCCUPIED
	} else {
		occupiedOrNot = traits.Occupancy_UNOCCUPIED
	}
	return &traits.Occupancy{
		State:       occupiedOrNot,
		PeopleCount: int32(n),
		Reasons:     []string{"A random occupancy value"},
	}
}
