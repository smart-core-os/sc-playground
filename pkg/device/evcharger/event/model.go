package event

import (
	"time"

	"github.com/smart-core-os/sc-playground/pkg/timeline"
)

// Model collects events that are relevant at a specific point in time.
type Model struct {
	TL  timeline.TL
	now time.Time

	PluggedIn    *PlugIn
	ChargeStared *ChargeStart
	ModeChanges  []*ModeChange // In descending-time order. ModeChanges[0] is the newest change
	Unplugged    *Unplug
}

func (m Model) At() time.Time {
	return m.now
}

func (m *Model) Scrub(t time.Time) error {
	// reset state
	m.Clear()
	m.now = t

	more := true
	for more {
		for _, o := range m.TL.At(t) {
			switch e := o.(type) {
			case PlugIn:
				m.PluggedIn = &e
				return nil // our model starts at the plugin event, OR
			case Unplug:
				m.Unplugged = &e // the Unplugged event.
				return nil
			case ChargeStart:
				if m.ChargeStared == nil {
					m.ChargeStared = &e
				}
			case ModeChange:
				m.ModeChanges = append(m.ModeChanges, &e)
			}
		}

		// keep looking
		t, more = m.TL.Previous(t)
	}

	return nil
}

func (m *Model) Clear() {
	m.PluggedIn = nil
	m.Unplugged = nil
	m.ChargeStared = nil
	m.ModeChanges = nil
}

func (m *Model) IsIdle() bool {
	return m.Unplugged != nil || m.PluggedIn == nil
}

func (m *Model) IsPluggedIn() bool {
	return m.PluggedIn != nil
}

func (m *Model) IsCharging() bool {
	return m.ChargeStared != nil
}
