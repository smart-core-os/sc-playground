package metadata

import (
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/metadata"
	"github.com/smart-core-os/sc-playground/pkg/node"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func Activate(n *node.Node) {
	r := metadata.NewApiRouter(
		metadata.WithMetadataApiClientFactory(func(name string) (traits.MetadataApiClient, error) {
			model := metadata.NewModel(resource.WithInitialValue(autoGeneratedDeviceMetadata(name)))
			return metadata.WrapApi(metadata.NewModelServer(model)), nil
		}),
		n.AnnounceOnRouterChange(trait.Metadata, node.NoAddMetadata()),
	)
	n.AddRouter(r)
	n.AddTraitFactory(trait.Metadata, func(name string, _ proto.Message) error {
		_, err := r.Get(name)
		return err
	})
	n.AddClientFactory(trait.Metadata, func(conn *grpc.ClientConn) interface{} {
		return traits.NewMetadataApiClient(conn)
	})
}

func autoGeneratedDeviceMetadata(name string) *traits.Metadata {
	return &traits.Metadata{
		Name: name,
		Appearance: &traits.Metadata_Appearance{
			Description: "Auto-generated device",
		},
		Traits: []*traits.TraitMetadata{
			{Name: string(trait.Metadata), More: node.AutoTraitMetadata},
		},
	}
}
