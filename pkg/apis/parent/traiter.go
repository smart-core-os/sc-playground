package parent

type Traiter interface {
	Trait(name string, traits ...string)
}
