package boot

import (
	"github.com/smart-core-os/sc-playground/pkg/device/evcharger"
	"github.com/smart-core-os/sc-playground/pkg/node"
	"github.com/smart-core-os/sc-playground/pkg/sim"
	"github.com/smart-core-os/sc-playground/pkg/sim/input"
	"github.com/smart-core-os/sc-playground/pkg/timeline/skiplisttl"
)

func CreateSimulation(n *node.Node) (*sim.Loop, error) {
	tl := skiplisttl.New()
	inputQueue := input.NewQueue()

	evcharger.New("sim/CHRG-01", tl, evcharger.WithInputDispatcher(inputQueue)).Publish(n)

	mainLoop := sim.NewLoop(tl, n, inputQueue)
	return mainLoop, nil
}
