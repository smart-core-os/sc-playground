package skiplisttl

import (
	"time"

	"github.com/MauriceGit/skiplist"
)

// Opt allows to configure how TL behaves.
// See With... functions for details.
type Opt func(tl *TL)

// DefaultOpts contains the default configuration for a TL.
// This keys time at a millisecond resolution.
var DefaultOpts = []Opt{
	WithSkipList(skiplist.New),
	WithKeyFunc(func(t time.Time) float64 {
		return float64(t.UnixMilli())
	}),
}

// WithKeyFunc configures the TL to use this function to convert time.Time into float64 keys for the skip list.
// The default is to use Time.UnixMilli().
func WithKeyFunc(keyFunc KeyFunc) Opt {
	return func(tl *TL) {
		tl.key = keyFunc
	}
}

// WithSkipList allows you to provide a customised version of the underlying skiplist.SkipList.
// Typically used to set a seed for testing.
// Defaults to skiplist.New().
func WithSkipList(sl func() skiplist.SkipList) Opt {
	return func(tl *TL) {
		underlying := sl()
		tl.items = &fixedSkipList{&underlying}
	}
}
