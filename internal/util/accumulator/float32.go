// Package accumulator implements helper types for iterative computations.
package accumulator

// Float32 implements a reduction operation over float32 using an arbitrary reduction function
type Float32 struct {
	Reduce func(accOld, x float32) (accNew float32)
	acc    float32
	init   bool
}

func (r *Float32) Accumulate(values ...float32) {
	for _, value := range values {
		if r.init {
			r.acc = r.Reduce(r.acc, value)
		} else {
			r.acc = value
			r.init = true
		}
	}
}

func (r *Float32) Get() (float32, bool) {
	return r.acc, r.init
}

func (r *Float32) GetOrDefault(def float32) float32 {
	v, ok := r.Get()
	if ok {
		return v
	} else {
		return def
	}
}
