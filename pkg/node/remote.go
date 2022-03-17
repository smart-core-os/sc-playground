package node

import (
	"context"
	"crypto/tls"

	"google.golang.org/grpc"
)

type remoteNode struct {
	endpoint string
	tls      *tls.Config
	conn     *grpc.ClientConn
	dialErr  error
}

type RemoteOption func(n *remoteNode)

func (n *Node) ResolveRemoteConn(ctx context.Context, endpoint string, opts ...RemoteOption) (*grpc.ClientConn, error) {
	// todo: check tls and update as needed
	rn := &remoteNode{endpoint: endpoint}
	for _, opt := range opts {
		opt(rn)
	}

	for _, node := range n.remoteNodes {
		if node.endpoint == endpoint {
			return node.conn, node.dialErr
		}
	}

	// create it
	n.remoteNodes = append(n.remoteNodes, createRemoteConnection(ctx, rn))
	return rn.conn, rn.dialErr
}

func createRemoteConnection(ctx context.Context, node *remoteNode) remoteNode {
	conn, err := grpc.DialContext(ctx, node.endpoint, grpc.WithInsecure())
	node.conn = conn
	node.dialErr = err
	return *node
}
