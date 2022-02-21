package evcharger

import (
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-playground/pkg/node"
)

func newMetadata(name string) *traits.Metadata {
	return &traits.Metadata{
		Name:         name,
		Traits:       metadataTraits(),
		Appearance:   metadataAppearance(),
		Location:     nil,
		Id:           nil,
		Product:      nil,
		Revision:     nil,
		Installation: nil,
		Nics:         nil,
		Membership:   nil,
		More:         metadataMore(),
	}
}

func metadataTraits() []*traits.TraitMetadata {
	return []*traits.TraitMetadata{
		{Name: string(trait.Electric)},
		{Name: string(trait.EnergyStorage)},
		{Name: string(trait.Metadata)},
	}
}

func metadataAppearance() *traits.Metadata_Appearance {
	return &traits.Metadata_Appearance{
		Title:       "EV Charger",
		Description: "Charges electric vehicles",
	}
}

func metadataMore() map[string]string {
	return map[string]string{
		node.MetadataDeviceType: "evcharger",
		node.MetadataRealism:    node.MetadataRealismVirtual,
	}
}
