package registry

// Slice is a Registry backed by a []Registry.
// Register calls Register for each Registry in the backing slice.
type Slice []Registry

func (r Slice) Register(name string, impl interface{}) (err error) {
	for _, reg := range r {
		err = reg.Register(name, impl)
	}
	return err
}

func (r *Slice) Add(registry Registry) {
	*r = append(*r, registry)
}
