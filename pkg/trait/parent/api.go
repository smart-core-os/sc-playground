package parent

import (
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/parent"
	"github.com/smart-core-os/sc-playground/pkg/node"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func Activate(n *node.Node) {
	r := parent.NewApiRouter(
		parent.WithParentApiClientFactory(func(name string) (traits.ParentApiClient, error) {
			return parent.WrapApi(parent.NewModelServer(parent.NewModel())), nil
		}),
		n.AnnounceOnRouterChange(trait.Parent),
	)
	n.AddRouter(r)
	n.AddTraitFactory(trait.Parent, func(name string, _ proto.Message) error {
		_, err := r.Get(name)
		return err
	})
	n.AddClientFactory(trait.Parent, func(conn *grpc.ClientConn) interface{} {
		return traits.NewParentApiClient(conn)
	})
}
