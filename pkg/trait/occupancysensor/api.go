package occupancysensor

import (
	"log"
	"math/rand"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/router"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"
	"github.com/smart-core-os/sc-golang/pkg/wrap"
	"github.com/smart-core-os/sc-playground/pkg/node"
	"google.golang.org/protobuf/proto"
)

func Activate(n *node.Node) {
	// handle random changes in occupancy
	type device struct {
		api  *occupancysensor.Model
		name string
	}
	var devices []device
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
			initial := randomOccupancy()
			log.Printf("Creating OccupancyApiClient(%v) %v (people=%v)", name, initial.State, initial.PeopleCount)
			model := occupancysensor.NewModel(initial)
			return occupancysensor.WrapApi(occupancysensor.NewModelServer(model)), nil
		}),
		router.WithOnCommit(func(name string, client interface{}) {
			model, ok := wrap.UnwrapFully(client).(*occupancysensor.Model)
			if !ok {
				return
			}

			occupancy, err := model.GetOccupancy()
			if err == nil {
				log.Printf("OccupancySensorApiClient(%v) auto-created %v (people=%v)", name, occupancy.State, occupancy.PeopleCount)
			} else {
				log.Printf("OccupancySensorApiClient(%v) auto-created (%v)", name, err)
			}
			devices = append(devices, device{api: model, name: name})
			n.Announce(name, node.HasTrait(trait.OccupancySensor))
		}),
	)
	n.AddRouter(r)
	n.AddTraitFactory(trait.OccupancySensor, func(name string, _ proto.Message) error {
		_, err := r.Get(name)
		return err
	})
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
