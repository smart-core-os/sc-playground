package playpb

import (
	"context"

	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-playground/pkg/node"
	"google.golang.org/grpc"
)

type Server struct {
	UnimplementedPlaygroundApiServer

	node *node.Node
}

func New(node *node.Node) *Server {
	return &Server{node: node}
}

func (s *Server) Register(server *grpc.Server) {
	RegisterPlaygroundApiServer(server, s)
}

func (s *Server) AddDeviceTrait(_ context.Context, request *AddDeviceTraitRequest) (*AddDeviceTraitResponse, error) {
	return &AddDeviceTraitResponse{}, s.node.CreateDeviceTrait(request.GetName(), trait.Name(request.GetTraitName()), nil)
}

func (s *Server) ListSupportedTraits(context.Context, *ListSupportedTraitsRequest) (*ListSupportedTraitsResponse, error) {
	res := &ListSupportedTraitsResponse{}
	for _, name := range s.node.SupportedDeviceTraits() {
		res.TraitName = append(res.TraitName, string(name))
	}
	return res, nil
}
