package node

import (
	"log"
	"time"

	"github.com/smart-core-os/sc-golang/pkg/router"
	"github.com/smart-core-os/sc-golang/pkg/server"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/parent"
	"github.com/smart-core-os/sc-playground/pkg/sim/scrub"
	"google.golang.org/grpc"
)

// Node represents a unit of control for a smart core server.
// Each node has collection of supported apis, represented by router.Router instances.
// When new devices are created they should call Announce and with the features relevant to the device.
type Node struct {
	name string

	children *parent.Model // lazy, initialised when addChildTrait or Register are called
	routers  []router.Router

	t         time.Time // last scrub time, so we can update added models
	scrubbers scrub.Slice
}

// New creates a new Node with the given device name.
func New(name string) *Node {
	return &Node{name: name}
}

// Register implements server.GrpcApi and registers all supported routers with s.
func (n *Node) Register(s *grpc.Server) {
	n.parent() // force the parent api to be initialised
	for _, r := range n.routers {
		if api, ok := r.(server.GrpcApi); ok {
			api.Register(s)
		}
	}
}

// AddRouter adds a router.Router to the published API of this node.
// AddRouter should not be called after Register is called.
func (n *Node) AddRouter(r ...router.Router) {
	n.routers = append(n.routers, r...)
}

// Announce adds a new device with the given features to this node.
// You may call Announce with the same name as a known device to add additional features, for example new traits.
func (n *Node) Announce(name string, features ...Feature) {
	a := &announcement{name: name}
	for _, feature := range features {
		feature.apply(a)
	}
	for _, t := range a.traits {
		log.Printf("%v now implements %v\n", name, t.name)

		if !t.noAddChildTrait && name != n.name {
			n.addChildTrait(a.name, t.name)
		}
		for _, client := range t.clients {
			n.addRoute(a.name, client)
		}
		if !t.noAddMetadata {
			md := t.metadata
			if md == nil {
				md = AutoTraitMetadata
			}
			if err := n.addTraitMetadata(name, t.name, md); err != nil {
				if err != MetadataTraitNotSupported {
					log.Printf("%v %v: %v", name, t.name, err)
				}
			}
		}
	}
	for _, simulator := range a.simulators {
		n.addScrubber(simulator)
	}
}

// Scrub implements the scrub.Scrubber interface and calls Scrub on each device with a simulation feature.
func (n *Node) Scrub(t time.Time) error {
	n.t = t
	return n.scrubbers.Scrub(t)
}

func (n *Node) addRoute(name string, impl interface{}) (added bool) {
	for _, r := range n.routers {
		if r.HoldsType(impl) {
			r.Add(name, impl)
			added = true
		}
	}
	return
}

func (n *Node) addChildTrait(name string, traitName ...trait.Name) {
	n.parent().AddChildTrait(name, traitName...)
}

func (n *Node) addScrubber(s scrub.Scrubber) {
	n.scrubbers = append(n.scrubbers, s)
	if !n.t.IsZero() {
		_ = s.Scrub(n.t)
	}
}

func (n *Node) parent() *parent.Model {
	if n.children == nil {
		// add this model as a device
		n.children = parent.NewModel()
		client := parent.WrapApi(parent.NewModelServer(n.children))
		n.Announce(n.name, HasTrait(trait.Parent, WithClients(client)))
	}
	return n.children
}
