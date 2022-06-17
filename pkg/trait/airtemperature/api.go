package airtemperature

import (
	"github.com/smart-core-os/sc-api/go/types"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/router"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/airtemperature"
	"github.com/smart-core-os/sc-golang/pkg/wrap"
	"github.com/smart-core-os/sc-playground/pkg/node"
	"google.golang.org/protobuf/proto"
)

func Activate(n *node.Node) {
	// handle random changes in occupancy
	type device struct {
		api  *airtemperature.Model
		name string
	}
	var devices []device
	go func() {
		for range time.Tick(30 * time.Second) {
			if len(devices) > 0 {
				i := rand.Intn(len(devices))
				o := randomTemperature()
				log.Printf("UpdateAirTemperature(%v) temp:%v, humidity:%v", devices[i].name, o.AmbientTemperature, o.AmbientHumidity)
				update, _ := devices[i].api.UpdateAirTemperature(o)
				log.Printf("updated %v", update)
			}
		}
	}()

	r := airtemperature.NewApiRouter(
		airtemperature.WithAirTemperatureApiClientFactory(func(name string) (traits.AirTemperatureApiClient, error) {
			log.Printf("Creating AirTemperatureApiClient(%v)", name)
			initial := randomTemperature()
			model := airtemperature.NewModel(initial)
			return airtemperature.WrapApi(airtemperature.NewModelServer(model)), nil
		}),
		router.WithOnChange(func(change router.Change) {
			if !change.Auto {
				return
			}
			name := change.Name
			model, ok := wrap.UnwrapFully(change.New).(*airtemperature.Model)
			if !ok {
				return
			}

			currentVal, err := model.GetAirTemperature()
			if err != nil {
				log.Printf("AirTemperatureApiClient(%v) auto-created (%v)", name, err)
			} else {
				log.Printf("AirTemperatureApiClient(%v) auto-created %v", name, currentVal.TemperatureGoal)
			}
			devices = append(devices, device{api: model, name: name})
			n.Announce(name, node.HasTrait(trait.AirTemperature))
		}),
	)
	n.AddRouter(r)
	n.AddTraitFactory(trait.AirTemperature, func(name string, _ proto.Message) error {
		_, err := r.Get(name)
		return err
	})
	n.AddClientFactory(trait.AirTemperature, func(conn *grpc.ClientConn) interface{} {
		return traits.NewAirTemperatureApiClient(conn)
	})
}

func randomTemperature() *traits.AirTemperature {
	var temp = 15 + rand.Float64()*10
	var hum float32 = rand.Float32()

	return &traits.AirTemperature{
		Mode:               traits.AirTemperature_AUTO,
		TemperatureGoal:    &traits.AirTemperature_TemperatureSetPoint{TemperatureSetPoint: &types.Temperature{ValueCelsius: 21.5}},
		AmbientTemperature: &types.Temperature{ValueCelsius: temp},
		AmbientHumidity:    &hum,
	}
}
