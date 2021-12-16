package broadcast

import (
	"github.com/smart-core-os/sc-golang/pkg/time/clock"
	"testing"
)

func TestVariable_Listen(t *testing.T) {
	v := NewVariable(clock.Real(), "foo")

	listner := v.Listen()
	v.Set("bar")

	change := <-listner.C

	if change.OldValue != "foo" || change.NewValue != "bar" {
		t.Error(change)
	}
}
