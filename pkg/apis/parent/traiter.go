package parent

import "github.com/smart-core-os/sc-golang/pkg/trait"

type Traiter interface {
	Trait(name string, traits ...trait.Name)
}
