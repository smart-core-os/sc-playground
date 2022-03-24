package playpb

import (
	"context"

	"github.com/smart-core-os/sc-golang/pkg/cmp"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-playground/pkg/node"
	"github.com/smart-core-os/sc-playground/pkg/sim/stats"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	UnimplementedPlaygroundApiServer

	node *node.Node

	performance *resource.Value // of *Performance
}

func New(node *node.Node) *Server {
	return &Server{
		node: node,
		performance: resource.NewValue(
			resource.WithInitialValue(new(Performance)),
			resource.WithMessageEquivalence(cmp.Equal(cmp.DurationValueWithinP(0.01))),
		),
	}
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

func (s *Server) AddRemoteDevice(ctx context.Context, req *AddRemoteDeviceRequest) (*AddRemoteDeviceResponse, error) {
	traitNames := req.TraitName
	if len(traitNames) == 0 {
		return nil, status.Errorf(codes.Unimplemented, "Trait discovery is not yet implements, please provide trait names explicitly")
	}

	var remoteOpts []node.RemoteOption
	if req.Insecure {
		remoteOpts = append(remoteOpts, node.WithRemoteInsecure())
	}
	if req.Tls != nil {
		if len(req.Tls.ServerCaCert) > 0 {
			remoteOpts = append(remoteOpts, node.WithRemoteServerCA([]byte(req.Tls.ServerCaCert)))
		}
	}
	conn, err := s.node.ResolveRemoteConn(ctx, req.Endpoint, remoteOpts...)
	if err != nil {
		return nil, err
	}

	var e error
	var features []node.Feature
	for _, traitName := range traitNames {
		client, err := s.node.CreateTraitClient(trait.Name(traitName), conn)
		if err != nil {
			e = multierr.Append(e, err)
			continue
		}
		features = append(features, node.HasTrait(trait.Name(traitName), node.WithClients(client), node.NoAddMetadata()))
	}
	s.node.Announce(req.Name, features...)
	return &AddRemoteDeviceResponse{}, e
}

func (s *Server) PullPerformance(_ *PullPerformanceRequest, stream PlaygroundApi_PullPerformanceServer) error {
	for change := range s.performance.Pull(stream.Context()) {
		performance := change.Value.(*Performance)
		err := stream.Send(&PullPerformanceResponse{
			Performance: performance,
			ChangeTime:  timestamppb.New(change.ChangeTime),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) Frame(f stats.FrameTiming) {
	_, _ = s.performance.Set(&Performance{
		Frame:   durationpb.New(f.Total),
		Capture: durationpb.New(f.Capture),
		Scrub:   durationpb.New(f.Scrub),
		Respond: durationpb.New(f.Respond),
		Idle:    durationpb.New(f.Idle),
	})
}
