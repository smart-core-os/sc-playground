package energystorage

import (
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait/energystorage"
	"github.com/tanema/gween/ease"
	"google.golang.org/protobuf/types/known/durationpb"
)

// Oscillator simulates an energy storage device by cycling between charging and discharging in an infinite loop.
type Oscillator struct {
	model *energystorage.Model

	cycleStart        time.Time
	chargeDuration    time.Duration
	fullDuration      time.Duration
	dischargeDuration time.Duration
	emptyDuration     time.Duration

	chargeRamp    ease.TweenFunc
	dischargeRamp ease.TweenFunc

	maxEnergyKWh          float32
	maxEnergyDistanceKm   float32
	unplugWhenNotCharging bool

	frameLength time.Duration
}

func NewOscillator(model *energystorage.Model, opts ...OscillatorOption) *Oscillator {
	o := &Oscillator{model: model}
	for _, opt := range DefaultOscillatorOptions {
		opt(o)
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (o *Oscillator) Scrub(t time.Time) error {
	// we loop, so normalise the time within our known universe
	totalLength := o.chargeDuration + o.fullDuration + o.dischargeDuration + o.emptyDuration
	playbackTime := t.Sub(o.cycleStart) % totalLength
	return o.scrub(playbackTime)
}

func (o *Oscillator) scrub(d time.Duration) (err error) {
	switch {
	case d < o.chargeDuration: // charging
		dSec := float32(d.Seconds())
		tSec := float32(o.chargeDuration.Seconds())

		_, err = o.model.UpdateEnergyLevel(&traits.EnergyLevel{
			Charging:  true,
			PluggedIn: true,
			Quantity: &traits.EnergyLevel_Quantity{
				Percentage: o.chargeRamp(dSec, 0, 100, tSec),
				EnergyKwh:  o.chargeRamp(dSec, 0, o.maxEnergyKWh, tSec),
			},
			UntilFull: &traits.EnergyLevel_UntilFull{
				Time:       durationpb.New(o.chargeDuration - d),
				DistanceKm: o.chargeRamp(dSec, 0, o.maxEnergyDistanceKm, tSec),
			},
		})
	case d < o.chargeDuration+o.fullDuration: // full
		_, err = o.model.UpdateEnergyLevel(&traits.EnergyLevel{
			Charging:  false,
			PluggedIn: !o.unplugWhenNotCharging,
			Quantity: &traits.EnergyLevel_Quantity{
				Percentage: 100,
				EnergyKwh:  o.maxEnergyKWh,
			},
			UntilFull: &traits.EnergyLevel_UntilFull{
				Time:       durationpb.New(0),
				DistanceKm: 0,
			},
		})
	case d < o.chargeDuration+o.fullDuration+o.dischargeDuration: // discharging
		durationDischarging := d - o.chargeDuration - o.fullDuration
		dSec := float32(durationDischarging.Seconds())
		tSec := float32(o.dischargeDuration.Seconds())

		_, err = o.model.UpdateEnergyLevel(&traits.EnergyLevel{
			Charging:  false,
			PluggedIn: !o.unplugWhenNotCharging,
			Quantity: &traits.EnergyLevel_Quantity{
				Percentage: 100 - o.dischargeRamp(dSec, 0, 100, tSec),
				EnergyKwh:  o.maxEnergyKWh - o.dischargeRamp(dSec, 0, o.maxEnergyKWh, tSec),
			},
			Remaining: &traits.EnergyLevel_Remaining{
				Time:       durationpb.New(o.dischargeDuration - durationDischarging),
				DistanceKm: o.maxEnergyDistanceKm - o.dischargeRamp(dSec, 0, o.maxEnergyDistanceKm, tSec),
			},
		})
	default: // empty
		_, err = o.model.UpdateEnergyLevel(&traits.EnergyLevel{
			Charging:  false,
			PluggedIn: !o.unplugWhenNotCharging,
			Quantity: &traits.EnergyLevel_Quantity{
				Percentage: 0,
				EnergyKwh:  0,
			},
			Remaining: &traits.EnergyLevel_Remaining{
				Time:       durationpb.New(0),
				DistanceKm: 0,
			},
		})
	}
	return err
}

type OscillatorOption func(o *Oscillator)

var DefaultOscillatorOptions = []OscillatorOption{
	WithCycleStartNow(),
	WithChargeDuration(60 * time.Second),
	WithFullDuration(20 * time.Second),
	WithDischargeDuration(4 * time.Minute),
	WithEmptyDuration(10 * time.Second),
	WithChargeRamp(ease.Linear),
	WithDischargeRamp(ease.Linear),
}

func WithCycleStart(cycleStart time.Time) OscillatorOption {
	return func(o *Oscillator) {
		o.cycleStart = cycleStart
	}
}

// WithCycleStartNow sets the cycleStart property to time.Now() each time an oscillator is created.
func WithCycleStartNow() OscillatorOption {
	return func(o *Oscillator) {
		o.cycleStart = time.Now()
	}
}
func WithChargeDuration(chargeDuration time.Duration) OscillatorOption {
	return func(o *Oscillator) {
		o.chargeDuration = chargeDuration
	}
}
func WithFullDuration(fullDuration time.Duration) OscillatorOption {
	return func(o *Oscillator) {
		o.fullDuration = fullDuration
	}
}
func WithDischargeDuration(dischargeDuration time.Duration) OscillatorOption {
	return func(o *Oscillator) {
		o.dischargeDuration = dischargeDuration
	}
}
func WithEmptyDuration(emptyDuration time.Duration) OscillatorOption {
	return func(o *Oscillator) {
		o.emptyDuration = emptyDuration
	}
}
func WithChargeRamp(chargeRamp ease.TweenFunc) OscillatorOption {
	return func(o *Oscillator) {
		o.chargeRamp = chargeRamp
	}
}
func WithDischargeRamp(dischargeRamp ease.TweenFunc) OscillatorOption {
	return func(o *Oscillator) {
		o.dischargeRamp = dischargeRamp
	}
}
func WithMaxEnergyKWh(maxEnergyKWh float32) OscillatorOption {
	return func(o *Oscillator) {
		o.maxEnergyKWh = maxEnergyKWh
	}
}
func WithMaxEnergyDistanceKm(maxEnergyDistanceKm float32) OscillatorOption {
	return func(o *Oscillator) {
		o.maxEnergyDistanceKm = maxEnergyDistanceKm
	}
}
func WithUnplugWhenNotCharging() OscillatorOption {
	return func(o *Oscillator) {
		o.unplugWhenNotCharging = true
	}
}
