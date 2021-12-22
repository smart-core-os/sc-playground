package broadcast

import (
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/olebedev/emitter"
	"github.com/smart-core-os/sc-golang/pkg/time/clock"
)

type MapEventType int

const (
	MapAddEvent MapEventType = iota
	MapRemoveEvent
	MapReplaceEvent
)

type MapEvent struct {
	Key       string
	Type      MapEventType
	Timestamp time.Time
	OldValue  interface{}
	NewValue  interface{}
}

// Map is a collection of key-value pairs, that can notify when an object is added, removed or replaced.
// Keys are strings; values can be any type.
type Map struct {
	m       sync.RWMutex
	members map[string]interface{}
	bus     *emitter.Emitter
	clock   clock.Clock
}

// NewMap constructs an empty set, using the given function to decide equality.
func NewMap(clk clock.Clock) *Map {
	s := &Map{
		members: make(map[string]interface{}),
		bus:     emitter.New(0),
		clock:   clk,
	}

	return s
}

// Set inserts or replaces an item in this Map.
// If there was no value with the same key in the map, then return value replaced will be false.
// If there was a member with the same key, the replaced will be true and the old value of that key will be returned.
func (s *Map) Set(key string, value interface{}) (old interface{}, replaced bool) {
	s.m.Lock()
	defer s.m.Unlock()

	existing, ok := s.members[key]
	s.members[key] = value

	if ok {
		s.bus.Emit(mapTopic, MapEvent{
			Key:       key,
			Type:      MapReplaceEvent,
			Timestamp: s.clock.Now(),
			OldValue:  existing,
			NewValue:  value,
		})
		return existing, true
	} else {
		s.bus.Emit(mapTopic, MapEvent{
			Key:       key,
			Type:      MapAddEvent,
			Timestamp: s.clock.Now(),
			NewValue:  value,
		})
		return nil, false
	}
}

func (s *Map) Remove(key string) (removed interface{}, ok bool) {
	s.m.Lock()
	defer s.m.Unlock()

	member, ok := s.members[key]
	if !ok {
		return nil, false
	}
	delete(s.members, key)
	s.bus.Emit(mapTopic, MapEvent{
		Key:       key,
		Type:      MapRemoveEvent,
		Timestamp: s.clock.Now(),
		OldValue:  member,
	})
	return member, true
}

// Members returns a copy of the current membership of the Map.
func (s *Map) Members() map[string]interface{} {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.copyMembers()
}

func (s *Map) copyMembers() map[string]interface{} {
	membersCopy := make(map[string]interface{})
	for k, v := range s.members {
		membersCopy[k] = v
	}

	return membersCopy
}

func (s *Map) GetMembersAndListen() (members map[string]interface{}, listener *MapListener) {
	s.m.RLock()
	defer s.m.RUnlock()

	members = s.copyMembers()
	listener = s.listen()

	return members, listener
}

func (s *Map) listen() *MapListener {
	dest := make(chan MapEvent)
	source := s.bus.On(mapTopic)
	done := make(chan struct{})

	go func() {
		defer close(dest)
		defer s.bus.Off(mapTopic, source)
		for {
			select {
			case <-done:
				return
			case event := <-source:
				dest <- event.Args[0].(MapEvent)
			}
		}
	}()

	listener := &MapListener{
		C:    dest,
		done: done,
	}
	runtime.SetFinalizer(listener, func(l *MapListener) {
		log.Println("logic error: broadcast.MapListener not closed before garbage collection")
		_ = l.Close()
	})
	return listener
}

func (s *Map) Listen() *MapListener {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.listen()
}

type MapListener struct {
	C         <-chan MapEvent
	closeOnce sync.Once
	done      chan<- struct{}
}

func (l *MapListener) Close() error {
	l.closeOnce.Do(func() {
		l.done <- struct{}{}
		close(l.done)
		runtime.SetFinalizer(l, nil)
	})
	return nil
}

const mapTopic = "mapChange"
