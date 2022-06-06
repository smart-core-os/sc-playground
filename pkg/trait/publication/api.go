package publication

import (
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/router"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/publication"
	"github.com/smart-core-os/sc-playground/pkg/node"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func Activate(n *node.Node) {
	r := publication.NewApiRouter(
		publication.WithPublicationApiClientFactory(func(name string) (traits.PublicationApiClient, error) {
			return publication.WrapApi(publication.NewModelServer(publication.NewModel())), nil
		}),
		router.WithOnChange(func(change router.Change) {
			if !change.Auto {
				return
			}
			n.Announce(change.Name, node.HasTrait(trait.Publication))
		}),
	)
	n.AddRouter(r)
	n.AddTraitFactory(trait.Publication, func(name string, _ proto.Message) error {
		_, err := r.Get(name)
		return err
	})
	n.AddClientFactory(trait.Publication, func(conn *grpc.ClientConn) interface{} {
		return traits.NewPublicationApiClient(conn)
	})
}
