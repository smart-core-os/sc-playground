package run

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io/fs"
	"os"

	"github.com/smart-core-os/sc-golang/pkg/server"
)

const (
	CertDir = "sc-playground/certs"
)

type Config struct {
	ctx context.Context

	defaultName string
	grpcApis    []server.GrpcApi

	grpcAddress  string
	httpAddress  string // includes http hosting and grpc-web over http
	httpsAddress string // includes https hosting and grpc-web over https

	hostedFS         fs.FS  // serve these files over http
	hostedFSNotFound []byte // serve these bytes on 404 for hosted fs
	httpHealthPath   string // expose this http path as a simple health api

	insecure bool // don't generate tls certs when they aren't provided
	mTLS     bool // turn on mutual-TLS support. Requires self-signed certs or the clientCaPEM option

	ca            *ca    // for generating certification
	clientCaPEM   []byte // CA cert used to verify client certs for mTLS
	grpcTlsConfig *tls.Config
	httpTlsConfig *tls.Config // If nil will use grpcTlsConfig instead

	certCacheDir string
	forceCertGen bool

	noHealth       bool // don't configure the health service
	noReflection   bool // don't configure the reflection service
	noGrpcWeb      bool // don't serve the grpc-web adapter
	noPlainGrpcWeb bool // don't serve grpc-web over plain http
}

func (c *Config) MarshalJSON() ([]byte, error) {
	type jt struct {
		GrpcAddress  string `json:"grpcAddress"`
		HttpAddress  string `json:"httpAddress"`
		HttpsAddress string `json:"httpsAddress"`
		SelfSigned   bool   `json:"selfSigned,omitempty"`
		Insecure     bool   `json:"insecure,omitempty"`
		MutualTLS    bool   `json:"mutualTls,omitempty"`
	}
	jo := jt{
		GrpcAddress:  c.grpcAddress,
		HttpAddress:  c.httpAddress,
		HttpsAddress: c.httpsAddress,
		SelfSigned:   c.ca != nil && !c.insecure,
		Insecure:     c.insecure,
		MutualTLS:    c.mTLS,
	}
	return json.Marshal(jo)
}

type ConfigOption func(*Config)

var DefaultOpts = []ConfigOption{
	WithContext(context.Background()),
	WithGrpcAddress(":9090"),
	WithHttpAddress(":8080"),
	WithHttpsAddress(":8443"),
}
var NilConfigOption ConfigOption = func(config *Config) {
	// do nothing
}

func WithContext(ctx context.Context) ConfigOption {
	return func(config *Config) {
		config.ctx = ctx
	}
}

// WithDefaultName configures the application with a default Smart Core name.
// The defaultName will be used as the default for requests that don't specify a name.
func WithDefaultName(defaultName string) ConfigOption {
	return func(config *Config) {
		config.defaultName = defaultName
	}
}

func WithApis(api ...server.GrpcApi) ConfigOption {
	return func(config *Config) {
		config.grpcApis = append(config.grpcApis, api...)
	}
}

func WithGrpcAddress(address string) ConfigOption {
	if address == "" {
		return NilConfigOption
	}
	return func(config *Config) {
		config.grpcAddress = address
	}
}

func WithHttpAddress(address string) ConfigOption {
	if address == "" {
		return NilConfigOption
	}
	return func(config *Config) {
		config.httpAddress = address
	}
}

func WithHttpsAddress(address string) ConfigOption {
	if address == "" {
		return NilConfigOption
	}
	return func(config *Config) {
		config.httpsAddress = address
	}
}

func WithHostedDir(dir string) ConfigOption {
	if dir == "" {
		return NilConfigOption
	}
	return WithHostedFS(os.DirFS(dir))
}

func WithHttpHealth(path string) ConfigOption {
	if path == "" {
		return NilConfigOption
	}
	return func(config *Config) {
		config.httpHealthPath = path
	}
}

func WithHostedFS(fs fs.FS) ConfigOption {
	return func(config *Config) {
		config.hostedFS = fs
	}
}

func WithHostedFSNotFound(html []byte) ConfigOption {
	return func(config *Config) {
		config.hostedFSNotFound = html
	}
}

func WithInsecure() ConfigOption {
	return func(config *Config) {
		config.insecure = true
	}
}

func WithGrpcTls(c *tls.Config) ConfigOption {
	return func(config *Config) {
		config.grpcTlsConfig = c
	}
}

func WithHttpTls(c *tls.Config) ConfigOption {
	return func(config *Config) {
		config.httpTlsConfig = c
	}
}

func WithForceCertGen() ConfigOption {
	return func(config *Config) {
		config.forceCertGen = true
	}
}

func WithCertCacheDir(certCacheDir string) ConfigOption {
	return func(config *Config) {
		config.certCacheDir = certCacheDir
	}
}

func WithMTLS() ConfigOption {
	return func(config *Config) {
		config.mTLS = true
	}
}

func NoHealth() ConfigOption {
	return func(config *Config) {
		config.noHealth = true
	}
}

func NoReflection() ConfigOption {
	return func(config *Config) {
		config.noReflection = true
	}
}

func NoGrpcWeb() ConfigOption {
	return func(config *Config) {
		config.noGrpcWeb = true
	}
}

func NoPlainGrpcWeb() ConfigOption {
	return func(config *Config) {
		config.noPlainGrpcWeb = true
	}
}
