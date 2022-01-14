package apis

type Traiter interface {
	Trait(name string, traits ...string)
}
