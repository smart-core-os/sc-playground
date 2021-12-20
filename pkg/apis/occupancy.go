package apis

import (
	"log"
	"math/rand"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/server"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancy"
)

func OccupancyApi() server.GrpcApi {
	// handle random changes in occupancy
	var devices []struct {
		api  *occupancy.Model
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

	r := occupancy.NewRouter()
	r.Factory = func(name string) (traits.OccupancySensorApiClient, error) {
		initial := randomOccupancy()
		log.Printf("Creating OccupancyApiClient(%v) %v (people=%v)", name, initial.State, initial.PeopleCount)
		api := occupancy.NewModel(initial)
		devices = append(devices, struct {
			api  *occupancy.Model
			name string
		}{api: api, name: name})
		return occupancy.Wrap(occupancy.NewModelServer(api)), nil
	}
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
