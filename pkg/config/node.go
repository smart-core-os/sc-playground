package config

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Node struct {
	Name      string     `json:"name,omitempty"`
	Address   string     `json:"address,omitempty"`
	Insecure  bool       `json:"insecure,omitempty"`
	TLSConfig *TLSConfig `json:"tls,omitempty"`
}

// Dial calls grpc.DialContext using options based on properties of n.
func (n Node) Dial(ctx context.Context) (*grpc.ClientConn, error) {
	if n.Address == "" {
		return nil, fmt.Errorf("node[%v].address is required", n.Name)
	}

	var opts []grpc.DialOption
	var err error

	opts, err = n.tlsDialOptions(ctx, opts)
	if err != nil {
		return nil, err
	}

	return grpc.DialContext(ctx, n.Address, opts...)
}

func (n Node) tlsDialOptions(ctx context.Context, opts []grpc.DialOption) ([]grpc.DialOption, error) {
	if n.Insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		tlsConfig := &tls.Config{}

		if n.TLSConfig != nil {
			serverCA, err := n.TLSConfig.ReadServerCACert(ctx)
			if err != nil {
				return opts, fmt.Errorf("node[%v].tls.serverCaCert %w", n.Name, err)
			}
			if serverCA != nil {
				certPool, err := x509.SystemCertPool()
				if err != nil {
					certPool = x509.NewCertPool()
				}
				if !certPool.AppendCertsFromPEM(serverCA) {
					return opts, fmt.Errorf("node[%v].tls.serverCaCert is not in PEM format or has no CERTIFICATE blocks", n.Name)
				}
				tlsConfig.RootCAs = certPool
			}
		}

		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	}
	return opts, nil
}

type TLSConfig struct {
	ServerCACert string `json:"serverCaCert,omitempty"`
}

func (t TLSConfig) ReadServerCACert(ctx context.Context) ([]byte, error) {
	if t.ServerCACert == "" {
		return nil, nil
	}

	// check if the cert is an embedded cert
	certBytes := []byte(t.ServerCACert)
	for {
		block, _ := pem.Decode(certBytes)
		if block == nil {
			break // the cert has no pem blocks
		}

		return certBytes, nil // rely on downstream processing of the cert to catch errors, like this not being a CERT block
	}

	// read from a http[s] url
	if strings.HasPrefix(t.ServerCACert, "http://") || strings.HasPrefix(t.ServerCACert, "https://") {
		// get from the network
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, t.ServerCACert, nil)
		if err != nil {
			return nil, err
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		if res.StatusCode < 200 || res.StatusCode >= 400 {
			return nil, fmt.Errorf("%v %v", res.StatusCode, res.Status)
		}
		return io.ReadAll(res.Body)
	}

	// read from a path
	return os.ReadFile(t.ServerCACert)
}

type NodeRef struct {
	Node
	NameOnly bool `json:"-"`
}

// Get resolves this reference using the given dict or it's own information if provided.
func (r NodeRef) Get(dict map[string]Node) (Node, bool) {
	if r.NameOnly {
		found, ok := dict[r.Name]
		return found, ok
	}
	return r.Node, true
}

// Normalize converts r into a NameOnly ref, placing inline data into dict.
func (r *NodeRef) Normalize(dict map[string]Node) (error, bool) {
	if r.NameOnly {
		return nil, false
	}
	name := r.Name
	if name == "" {
		return errors.New("no name"), false
	}
	_, ok := dict[name]
	if ok {
		return fmt.Errorf("duplicate %v", name), false
	}
	dict[name] = r.Node
	r.Node = Node{Name: name}
	return nil, true
}

func (r *NodeRef) UnmarshalJSON(bytes []byte) error {
	var name string
	if err := json.Unmarshal(bytes, &name); err == nil {
		r.Name = name
		r.NameOnly = true
		return nil
	}
	return json.Unmarshal(bytes, &r.Node)
}

func (r *NodeRef) MarshalJSON() ([]byte, error) {
	if r.NameOnly {
		return json.Marshal(r.Name)
	}
	return json.Marshal(r.Node)
}
