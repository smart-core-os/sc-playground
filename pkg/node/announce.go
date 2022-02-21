package node

import (
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-playground/pkg/sim/scrub"
)

// Announcer defines the Announce method.
// Calling Announce signals that a given named device has the collection of features provided.
type Announcer interface {
	// Announce signals that the named device has the given features.
	Announce(name string, features ...Feature)
}

type announcement struct {
	name       string
	simulators []scrub.Scrubber
	traits     []traitFeature
	clients    []interface{}
}

type traitFeature struct {
	name     trait.Name
	clients  []interface{}
	metadata map[string]string

	noAddChildTrait bool
	noAddMetadata   bool
}

// Feature describes some aspect of a named device.
type Feature interface{ apply(a *announcement) }

type featureFunc func(a *announcement)

func (f featureFunc) apply(a *announcement) {
	f(a)
}

// HasSimulation indicates that the device can be scrubbed using the given scrub.Scrubber.
func HasSimulation(s scrub.Scrubber) Feature {
	return featureFunc(func(a *announcement) {
		a.simulators = append(a.simulators, s)
	})
}

// HasClient indicates that the device has additional apis that aren't tied directly to a trait.
func HasClient(clients ...interface{}) Feature {
	return featureFunc(func(a *announcement) {
		a.clients = append(a.clients, clients...)
	})
}

// HasTrait indicates that the device implements the named trait.
func HasTrait(name trait.Name, opt ...TraitOption) Feature {
	return featureFunc(func(a *announcement) {
		feature := traitFeature{name: name}
		for _, option := range opt {
			option(&feature)
		}
		a.traits = append(a.traits, feature)
	})
}

// TraitOption controls how a Node behaves when presented with a new device trait.
type TraitOption func(t *traitFeature)

// WithClients indicates that the trait is implemented by these client instances.
// The clients will be added to the relevant routers when the trait is announced.
func WithClients(client ...interface{}) TraitOption {
	return func(t *traitFeature) {
		t.clients = append(t.clients, client...)
	}
}

// NoAddChildTrait instructs the Node not to add the trait to the nodes parent.Model.
func NoAddChildTrait() TraitOption {
	return func(t *traitFeature) {
		t.noAddChildTrait = true
	}
}

// NoAddMetadata instructs the Node not to add the trait to the nodes traits.Metadata.
func NoAddMetadata() TraitOption {
	return func(t *traitFeature) {
		t.noAddMetadata = true
	}
}

// WithTraitMetadata instructs the Node to use the given metadata when adding the trait to the nodes traits.Metadata.
// Metadata maps will be merged together, with conflicting keys in later calls overriding existing keys.
func WithTraitMetadata(md map[string]string) TraitOption {
	return func(t *traitFeature) {
		if t.metadata == nil {
			t.metadata = make(map[string]string)
		}
		for k, v := range md {
			t.metadata[k] = v
		}
	}
}
