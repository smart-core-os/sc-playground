package electric

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/time/clock"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/electric"
	"github.com/smart-core-os/sc-playground/pkg/node"
	"go.uber.org/multierr"
)

// Group presents a group of named ElectricApiClients as a unified device.
type Group struct {
	client traits.ElectricApiClient

	children []reducer

	closeCtx context.Context
	close    func()

	signals chan signal

	startOnce sync.Once
}

func NewGroup(client traits.ElectricApiClient) *Group {
	g := &Group{
		client:  client,
		signals: make(chan signal),
	}
	g.closeCtx, g.close = context.WithCancel(context.Background())
	g.children = []reducer{
		{"max", electric.NewModel(clock.Real()), max},
		{"min", electric.NewModel(clock.Real()), min},
		{"sum", electric.NewModel(clock.Real()), sum},
		{"average", electric.NewModel(clock.Real()), average},
	}

	return g
}

// Announce registers this group with the given node and router.
// It registers using the given name, as well as a number of child devices representing different combinations of the
// group members.
func (g *Group) Announce(name string, n *node.Node, r *electric.ApiRouter) {
	for i, child := range g.children {
		fqName := name + "/" + child.name
		client := electric.WrapApi(electric.NewModelServer(child.model))
		r.Add(fqName, client)
		n.Announce(fqName, node.HasTrait(trait.Electric, node.WithTraitMetadata(map[string]string{node.MetadataRealism: "computed"})))
		if i == 0 {
			// default child
			// todo: the default child should be a parent with g.children as children.
			r.Add(name, client)
			n.Announce(name, node.HasTrait(trait.Electric, node.WithTraitMetadata(map[string]string{node.MetadataRealism: "alias:./" + child.name})))
		}
	}
}

func (g *Group) Close() error {
	g.close()
	return nil
}

func (g *Group) Start() {
	g.startOnce.Do(func() {
		go g.run()
	})
}

func (g *Group) Add(name string) {
	g.Start()
	select {
	case <-g.closeCtx.Done():
	case g.signals <- memberAdded{name}:
	}
}

func (g *Group) Remove(name string) {
	g.Start()
	select {
	case <-g.closeCtx.Done():
	case g.signals <- memberRemoved{name}:
	}
}

type value struct {
	current float32
	err     error
	cancel  context.CancelFunc
}

func (g *Group) run() {
	memberCurrent := make(map[string]*value)
	noValue := errors.New("no value set yet")

	for {
		select {
		case <-g.closeCtx.Done():
			return
		case signal := <-g.signals:
			switch sig := signal.(type) {
			case memberAdded:
				name := sig.name
				childCtx, cancelStream := context.WithCancel(g.closeCtx)
				memberCurrent[name] = &value{err: noValue, cancel: cancelStream}
				go func() {
					res, err := g.client.GetDemand(childCtx, &traits.GetDemandRequest{Name: name})
					g.recordValueChange(name, res, err)

					stream, err := g.client.PullDemand(childCtx, &traits.PullDemandRequest{Name: name})
					if err != nil {
						g.recordValueChange(name, nil, err)
						return
					}
					for {
						recv, err := stream.Recv()
						if err != nil {
							g.recordValueChange(name, nil, err)
							return
						}
						latestChange := recv.Changes[len(recv.Changes)-1]
						g.recordValueChange(name, latestChange.Demand, nil)
					}
				}()
			case memberRemoved:
				if old, ok := memberCurrent[sig.name]; ok {
					old.cancel()
					delete(memberCurrent, sig.name)
					if err := g.notifyChildren(memberCurrent); err != nil {
						log.Printf("Error notifying children of member removal %v", err)
					}
				}
			case demandCurrentChange:
				val, ok := memberCurrent[sig.name]
				if !ok {
					continue // member was removed before we could record the change
				}
				val.current = sig.current
				val.err = sig.err
				if err := g.notifyChildren(memberCurrent); err != nil {
					log.Printf("Error notifying children of current change %v", err)
				}
			}
		}
	}
}

func (g *Group) recordValueChange(name string, demand *traits.ElectricDemand, err error) {
	select {
	case g.signals <- demandCurrentChange{name: name, current: demand.GetCurrent(), err: err}:
	case <-g.closeCtx.Done():
	}
}

func (g *Group) notifyChildren(values map[string]*value) error {
	// collect all the non-error values together
	var currents []float32
	for _, v := range values {
		if v.err != nil {
			continue
		}
		currents = append(currents, v.current)
	}

	var err error
	for _, child := range g.children {
		if err2 := child.Publish(currents); err2 != nil {
			err = multierr.Append(err, err2)
		}
	}
	return err
}

type reducer struct {
	name  string
	model *electric.Model
	fn    func(readings []float32) float32
}

func (r reducer) Publish(readings []float32) error {
	reading := r.fn(readings)
	_, err := r.model.UpdateDemand(&traits.ElectricDemand{Current: reading}, resource.WithUpdatePaths("current"))
	return err
}

func max(readings []float32) float32 {
	var max float32
	for _, reading := range readings {
		if reading > max {
			max = reading
		}
	}
	return max
}

func min(readings []float32) float32 {
	var min float32
	var minSet bool
	for _, reading := range readings {
		if reading == 0 {
			continue // don't count devices without load
		}
		if !minSet || reading < min {
			minSet = true
			min = reading
		}
	}
	return min
}

func sum(readings []float32) float32 {
	var total float32
	for _, reading := range readings {
		total += reading
	}
	return total
}

// average sums the non-zero readings and divides by how many were summed.
func average(readings []float32) float32 {
	var nonZeroReadings float32
	var total float32
	for _, reading := range readings {
		if reading != 0 {
			nonZeroReadings++
			total += reading
		}
	}
	return total / nonZeroReadings
}

type signal interface {
	implementsSignal()
}

type demandCurrentChange struct {
	name    string
	current float32
	err     error
}

func (s demandCurrentChange) implementsSignal() {
}

type memberRemoved struct {
	name string
}

func (s memberRemoved) implementsSignal() {
}

type memberAdded struct {
	name string
}

func (s memberAdded) implementsSignal() {
}
