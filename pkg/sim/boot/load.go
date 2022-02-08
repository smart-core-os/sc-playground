package boot

import (
	"github.com/smart-core-os/sc-playground/pkg/apis/registry"
	"github.com/smart-core-os/sc-playground/pkg/device/evcharger"
	"github.com/smart-core-os/sc-playground/pkg/sim"
	"github.com/smart-core-os/sc-playground/pkg/sim/input"
	"github.com/smart-core-os/sc-playground/pkg/timeline/skiplisttl"
)

func CreateSimulation(reg registry.Registry) (*sim.Loop, error) {
	tl := skiplisttl.New()
	inputQueue := input.NewQueue()

	chrg01 := evcharger.New("sim/CHRG-01", tl, evcharger.WithInputDispatcher(inputQueue))

	if err := chrg01.Publish(reg); err != nil {
		return nil, err
	}

	mainLoop := sim.NewLoop(tl, chrg01, inputQueue)
	return mainLoop, nil
}
