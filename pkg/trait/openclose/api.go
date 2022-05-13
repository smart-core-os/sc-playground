package openclose

import (
	"github.com/smart-core-os/sc-golang/pkg/trait/openclose"
	"github.com/smart-core-os/sc-playground/pkg/node"
)

func Activate(n *node.Node) {
	n.AddRouter(openclose.NewApiRouter())
}
