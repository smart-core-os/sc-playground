package broadcast

import (
	"github.com/olebedev/emitter"
	"github.com/smart-core-os/sc-playground/internal/util/clock"
	"log"
	"runtime"
	"sync"
	"time"
)

const variableTopic = "variable"

type VariableEvent struct {
	NewValue  interface{}
	OldValue  interface{}
	Timestamp time.Time
}

type Variable struct {
	m     sync.RWMutex
	em    *emitter.Emitter
	clock clock.Clock

	value   interface{}
	changed time.Time
}

func NewVariable(c clock.Clock, initialValue interface{}) *Variable {
	v := &Variable{
		em:      emitter.New(0),
		clock:   c,
		value:   initialValue,
		changed: c.Now(),
	}

	return v
}

func (v *Variable) Set(value interface{}) {
	v.m.Lock()
	defer v.m.Unlock()

	v.changed = v.clock.Now()
	v.em.Emit(variableTopic, VariableEvent{
		NewValue:  value,
		OldValue:  v.value,
		Timestamp: v.changed,
	})
	v.value = value
}

func (v *Variable) Value() interface{} {
	value, _ := v.Get()
	return value
}

func (v *Variable) Get() (value interface{}, changed time.Time) {
	v.m.RLock()
	defer v.m.RUnlock()

	return v.value, v.changed
}

func (v *Variable) Listen() *Listener {
	_, listener := v.GetAndListen()
	return listener
}

func (v *Variable) GetAndListen() (value interface{}, listener *Listener) {
	v.m.RLock()
	defer v.m.RUnlock()

	dest := make(chan VariableEvent)
	source := v.em.On(variableTopic)
	done := make(chan struct{})

	// wrapper goroutine to convert event types
	go func() {
		defer close(dest)
		for {
			select {
			case <-done:
				v.em.Off(variableTopic, source)
				return
			case event := <-source:
				dest <- event.Args[0].(VariableEvent)
			}
		}
	}()

	listener = &Listener{
		done: done,
		C:    dest,
	}

	runtime.SetFinalizer(listener, func(l *Listener) {
		log.Println("logic error: broadcast.Listener was not closed before garbage collection")
		_ = l.Close()
	})

	return v.value, listener
}

type Listener struct {
	closeOnce sync.Once
	done      chan struct{}
	C         <-chan VariableEvent
}

// Close will free the listener. No more events will be sent on the channel.
// The returned error is always nil - signature is for compatibility with io.Closer.
func (l *Listener) Close() error {
	l.closeOnce.Do(func() {
		runtime.SetFinalizer(l, nil)
		l.done <- struct{}{}
		close(l.done)
	})
	return nil
}
