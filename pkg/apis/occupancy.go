package apis

import (
	"context"
	"log"
	"math/rand"
	"time"

	"git.vanti.co.uk/smartcore/sc-api/go/traits"
	"git.vanti.co.uk/smartcore/sc-golang/pkg/memory"
	"git.vanti.co.uk/smartcore/sc-golang/pkg/router"
	"git.vanti.co.uk/smartcore/sc-golang/pkg/server"
	"git.vanti.co.uk/smartcore/sc-golang/pkg/wrap"
)

func OccupancyApi() server.GrpcApi {
	// handle random changes in occupancy
	var devices []struct {
		api  *memory.OccupancySensorApi
		name string
	}
	go func() {
		for range time.Tick(30 * time.Second) {
			if len(devices) > 0 {
				i := rand.Intn(len(devices))
				o := randomOccupancy()
				log.Printf("SetOccupancy(%v) %v (people=%v)", devices[i].name, o.State, o.PeopleCount)
				devices[i].api.SetOccupancy(context.Background(), o)
			}
		}
	}()

	r := router.NewOccupancySensorApiRouter()
	r.Factory = func(name string) (traits.OccupancySensorApiClient, error) {
		occupancy := randomOccupancy()
		log.Printf("Creating OccupancyApiClient(%v) %v (people=%v)", name, occupancy.State, occupancy.PeopleCount)
		api := memory.NewOccupancyApi(occupancy)
		devices = append(devices, struct {
			api  *memory.OccupancySensorApi
			name string
		}{api: api, name: name})
		return wrap.OccupancySensorApiServer(api), nil
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
