package node

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type remoteNode struct {
	endpoint string
	tls      *tls.Config
	conn     *grpc.ClientConn
	dialErr  error
}

type RemoteOption func(n *remoteNode)

func WithRemoteServerCA(ca []byte) RemoteOption {
	return func(n *remoteNode) {
		if n.tls == nil {
			n.tls = &tls.Config{}
		}
		rootCAs, err := x509.SystemCertPool()
		if err != nil {
			log.Printf("Error reading system cert pool %v", err)
			rootCAs = x509.NewCertPool()
		}
		if !rootCAs.AppendCertsFromPEM(ca) {
			log.Printf("Unable to append ca to system certs")
		}
		n.tls.RootCAs = rootCAs
	}
}

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
	var grpcOpts []grpc.DialOption
	if node.tls == nil {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
	} else {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(node.tls)))
	}
	conn, err := grpc.DialContext(ctx, node.endpoint, grpcOpts...)
	node.conn = conn
	node.dialErr = err
	return *node
}
