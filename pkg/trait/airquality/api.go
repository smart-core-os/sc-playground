package airquality

import (
	"github.com/smart-core-os/sc-golang/pkg/trait/airqualitysensor"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/router"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/wrap"
	"github.com/smart-core-os/sc-playground/pkg/node"
	"google.golang.org/protobuf/proto"
)

func Activate(n *node.Node) {
	type device struct {
		api  *airqualitysensor.Model
		name string
	}
	var devices []device
	go func() {
		for range time.Tick(30 * time.Second) {
			if len(devices) > 0 {
				i := rand.Intn(len(devices))
				o := randomAirQuality()
				log.Printf("UpdateAirQuality(%v) co2:%v, voc:%v", devices[i].name, o.CarbonDioxideLevel, o.VolatileOrganicCompounds)
				update, _ := devices[i].api.UpdateAirQuality(o)
				log.Printf("updated %v", update)
			}
		}
	}()

	r := airqualitysensor.NewApiRouter(
		airqualitysensor.WithAirQualitySensorApiClientFactory(func(name string) (traits.AirQualitySensorApiClient, error) {
			log.Printf("Creating AirQualitySensorApiClient(%v)", name)
			initial := randomAirQuality()
			model := airqualitysensor.NewModel(initial)
			return airqualitysensor.WrapApi(airqualitysensor.NewModelServer(model)), nil
		}),
		router.WithOnChange(func(change router.Change) {
			if !change.Auto {
				return
			}
			name := change.Name
			model, ok := wrap.UnwrapFully(change.New).(*airqualitysensor.Model)
			if !ok {
				return
			}

			currentVal, err := model.GetAirQuality()
			if err != nil {
				log.Printf("AirQualitySensorApiClient(%v) auto-created (%v)", name, err)
			} else {
				log.Printf("AirQualitySensorApiClient(%v) auto-created %v", name, currentVal.AirPressure)
			}
			devices = append(devices, device{api: model, name: name})
			n.Announce(name, node.HasTrait(trait.AirQualitySensor))
		}),
	)
	n.AddRouter(r)
	n.AddTraitFactory(trait.AirQualitySensor, func(name string, _ proto.Message) error {
		_, err := r.Get(name)
		return err
	})
	n.AddClientFactory(trait.AirQualitySensor, func(conn *grpc.ClientConn) interface{} {
		return traits.NewAirQualitySensorApiClient(conn)
	})
}

func randomAirQuality() *traits.AirQuality {
	// < 1000 is OK, 400-1000 is normal inside
	var co2 = 600 + (rand.Float32()-0.5)*200 // 400-800
	// couldn't find much on what is OK, but <0.5 I think
	var voc = rand.Float32() * 0.5
	// 1013 is atmospheric pressure at sea level
	var pressure = 1013 + (rand.Float32()-0.5)*10 // 1003-1023

	return &traits.AirQuality{
		AirPressure:              &pressure,
		CarbonDioxideLevel:       &co2,
		VolatileOrganicCompounds: &voc,
	}
}
