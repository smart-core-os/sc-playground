package sim

import (
	"os"

	"github.com/smart-core-os/sc-playground/pkg/sim/stats"
)

type Option func(l *Loop)

var DefaultOptions = []Option{
	WithFramer(stats.CFramer(60*5, stats.FFramer(os.Stdout))),
}

func WithFramer(f stats.Framer) Option {
	return func(l *Loop) {
		l.framer = f
	}
}

func applyOpts(l *Loop, opts ...Option) {
	for _, option := range DefaultOptions {
		option(l)
	}
	for _, opt := range opts {
		opt(l)
	}
}
