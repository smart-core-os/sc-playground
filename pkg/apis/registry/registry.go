package registry

type Publisher interface {
	Publish(reg Registry) error
}

type Registry interface {
	Register(name string, impl interface{}) error
}

type Adder interface {
	Add(registry Registry)
}
