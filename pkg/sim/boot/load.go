package boot

import (
	"fmt"

	"github.com/smart-core-os/sc-playground/pkg/device/evcharger"
	"github.com/smart-core-os/sc-playground/pkg/node"
	"github.com/smart-core-os/sc-playground/pkg/sim"
	"github.com/smart-core-os/sc-playground/pkg/sim/input"
	"github.com/smart-core-os/sc-playground/pkg/timeline/skiplisttl"
)

func CreateSimulation(n *node.Node, opts ...sim.Option) (*sim.Loop, error) {
	tl := skiplisttl.New()
	inputQueue := input.NewQueue()

	chargerCount := 0
	for i := 0; i < chargerCount; i++ {
		evcharger.New(fmt.Sprintf("sim/EVC-%02d", i+1), tl, evcharger.WithInputDispatcher(inputQueue)).Publish(n)
	}

	mainLoop := sim.NewLoop(tl, n, inputQueue, opts...)
	return mainLoop, nil
}
