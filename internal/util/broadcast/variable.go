package broadcast

import (
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/olebedev/emitter"
	"github.com/smart-core-os/sc-golang/pkg/time/clock"
)

const variableTopic = "variable"

type VariableEvent struct {
	NewValue  interface{}
	OldValue  interface{}
	Timestamp time.Time
}

// Variable holds an interface{} value. Whenever the value is changed, the Variable sends a VariableEvent
// to all its Listeners.
// Note that only method calls will actually trigger events. Do not mutate values that are stored in the variable,
// as this will not be thread safe without external synchronisation and will not send a VariableEvent.
// When you need to modify the variable value, call Set.
type Variable struct {
	m     sync.RWMutex
	em    *emitter.Emitter
	clock clock.Clock

	value   interface{}
	changed time.Time
}

// NewVariable constructs a Variable.
// The Variable initially stores the value given by initialValue
func NewVariable(c clock.Clock, initialValue interface{}) *Variable {
	v := &Variable{
		em:      emitter.New(0),
		clock:   c,
		value:   initialValue,
		changed: c.Now(),
	}

	return v
}

// Set will store the value in this Variable, and broadcast the change as a VariableEvent to all Listeners open
// on this variable.
// Set is thread-safe.
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

// Value returns the current value of the Variable.
func (v *Variable) Value() interface{} {
	value, _ := v.Get()
	return value
}

// Get returns the current value of the Variable, and the timestamp of last modification.
func (v *Variable) Get() (value interface{}, changed time.Time) {
	v.m.RLock()
	defer v.m.RUnlock()

	return v.value, v.changed
}

// Listen opens a listener on the Variable which will receive events describing any future changes to
// the Variable's value.
func (v *Variable) Listen() *VariableListener {
	_, listener := v.GetAndListen()
	return listener
}

// GetAndListen will return the current value of the Variable, and a listener that will receive events describing
// future changes. This operation happens atomically, so it is guaranteed that you will not miss any changes
// that happen to the value in between getting the value and creating the VariableListener.
func (v *Variable) GetAndListen() (value interface{}, listener *VariableListener) {
	v.m.RLock()
	defer v.m.RUnlock()

	dest := make(chan VariableEvent)
	source := v.em.On(variableTopic)
	done := make(chan struct{})

	// wrapper goroutine to convert event types
	go func() {
		defer close(dest)
		defer v.em.Off(variableTopic, source)
		for {
			select {
			case <-done:
				return
			case event := <-source:
				dest <- event.Args[0].(VariableEvent)
			}
		}
	}()

	listener = &VariableListener{
		done: done,
		C:    dest,
	}

	runtime.SetFinalizer(listener, func(l *VariableListener) {
		log.Println("logic error: broadcast.VariableListener was not closed before garbage collection")
		_ = l.Close()
	})

	return v.value, listener
}

// VariableListener allows receiving notification of all changes to an associated Variable.
// To obtain a VariableListener, call Variable.Listen or Variable.GetAndListen.
// It is important to free the VariableListener by calling Close after it is no longer needed.
type VariableListener struct {
	closeOnce sync.Once
	done      chan struct{}
	C         <-chan VariableEvent
}

// Close will free the listener. No more events will be sent on the channel.
// The returned error is always nil - signature is for compatibility with io.Closer.
func (l *VariableListener) Close() error {
	l.closeOnce.Do(func() {
		runtime.SetFinalizer(l, nil)
		// send through channel to create a synchronisation point
		l.done <- struct{}{}
		close(l.done)
	})
	return nil
}
