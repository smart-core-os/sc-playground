package node

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type remoteNode struct {
	endpoint string
	tls      *tls.Config
	insecure bool
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

func WithRemoteInsecure() RemoteOption {
	return func(n *remoteNode) {
		n.insecure = true
	}
}

func WithRemoteSkipVerify() RemoteOption {
	return func(n *remoteNode) {
		n.tls = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
}

func (n *Node) ResolveRemoteConn(ctx context.Context, endpoint string, opts ...RemoteOption) (*grpc.ClientConn, error) {
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
	if node.insecure {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		tlsConfig := node.tls
		if tlsConfig == nil {
			tlsConfig = &tls.Config{} // uses defaults, like system CA pool
		}
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	}
	conn, err := grpc.DialContext(ctx, node.endpoint, grpcOpts...)
	node.conn = conn
	node.dialErr = err
	if err == nil {
		go func() {
			// watch state changes
			var state connectivity.State
			for conn.WaitForStateChange(context.Background(), state) {
				state = conn.GetState()
				log.Printf("%v is %v", node.endpoint, state)
			}
		}()
	}
	return *node
}
