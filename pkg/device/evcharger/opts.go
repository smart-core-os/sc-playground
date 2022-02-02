package evcharger

var (
	DefaultRating  float32 = 60  // Amps
	DefaultVoltage float32 = 230 // Volts
)

type Opt func(*opt)
type opt struct {
	rating      float32
	voltage     float32
	dedicatedTL bool
}

var DefaultOpts = []Opt{
	WithVoltage(DefaultVoltage),
	WithRating(DefaultRating),
}

func WithRating(r float32) Opt {
	return func(o *opt) {
		o.rating = r
	}
}

func WithVoltage(v float32) Opt {
	return func(o *opt) {
		o.voltage = v
	}
}

func WithDedicatedTL() Opt {
	return func(o *opt) {
		o.dedicatedTL = true
	}
}
