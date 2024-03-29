package run

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/rs/cors"
	"github.com/smart-core-os/sc-golang/pkg/middleware/name"
	"github.com/smart-core-os/sc-golang/pkg/server"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func Serve(opts ...ConfigOption) error {
	config := new(Config)
	for _, opt := range DefaultOpts {
		opt(config)
	}
	for _, opt := range opts {
		opt(config)
	}

	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	done := make(chan error, 1)

	// messages we print at the end to say what just happened
	// keyed by address, value is a slice of the protocols served at that address
	runMsg := make(map[net.Addr][]string)
	addRunMsg := func(address net.Addr, protocol ...string) {
		runMsg[address] = append(runMsg[address], protocol...)
	}

	grpcLis, err := net.Listen("tcp", config.grpcAddress)
	if err != nil {
		return fmt.Errorf("grpc: %w", err)
	}
	var grpcOpts []grpc.ServerOption
	if config.defaultName != "" {
		grpcOpts = append(grpcOpts,
			grpc.ChainUnaryInterceptor(name.IfAbsentUnaryInterceptor(config.defaultName)),
			grpc.ChainStreamInterceptor(name.IfAbsentStreamInterceptor(config.defaultName)),
		)
	}
	if config.grpcTlsConfig != nil {
		addRunMsg(grpcLis.Addr(), "Secure gRPC")
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(config.grpcTlsConfig)))
	} else if !config.insecure {
		ca, serverCert, err := loadOrCreateCerts(config.forceCertGen, config.certCacheDir)
		if err != nil {
			return err
		}
		config.ca = ca
		config.grpcTlsConfig = &tls.Config{
			Certificates: []tls.Certificate{
				*serverCert,
			},
		}

		// mTLS
		if config.mTLS {
			config.grpcTlsConfig.ClientCAs, err = config.ca.Pool()
			config.grpcTlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
			addRunMsg(grpcLis.Addr(), "Self-signed mTLS gRPC")
		} else {
			addRunMsg(grpcLis.Addr(), "Self-signed gRPC")
		}
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(config.grpcTlsConfig)))
	} else {
		addRunMsg(grpcLis.Addr(), "Insecure gRPC")
	}

	// add some request debugging
	grpcOpts = append(grpcOpts, grpc.ChainStreamInterceptor(LogRepeatChanges()))

	grpcServer := grpc.NewServer(grpcOpts...)
	defer grpcServer.Stop()

	// register client apis
	for _, api := range config.grpcApis {
		api.Register(grpcServer)
	}

	// register infrastructure apis
	if !config.noReflection {
		reflection.Register(grpcServer)
	}
	// register health last so we can add all services
	if !config.noHealth {
		registerHealth(grpcServer)
	}

	go func() { done <- grpcServer.Serve(grpcLis) }()

	var webPageUrl *url.URL
	noHttp := config.noGrpcWeb && config.hostedFS == nil && config.httpHealthPath == ""
	if !noHttp {
		httpServer := http.Server{}
		defer httpServer.Close()

		var interceptors []HttpInterceptor
		// grpc-web handler
		if !config.noGrpcWeb {
			grpcWebServer := grpcweb.WrapServer(grpcServer, grpcweb.WithOriginFunc(func(origin string) bool {
				return true
			}))
			interceptors = append(interceptors, grpcWebInterceptor(grpcWebServer))
		}

		// use a ServeMux to allow for future (non-hosted) http endpoints
		mux := http.NewServeMux()
		// provided APIs used by client libraries to configure themselves
		mux.Handle("/__/playground/config.json", configJsonHandler(config))
		mux.Handle("/__/playground/ca-cert.pem", caCertHandler(config))
		mux.Handle("/__/playground/client.pem", clientCertHandler(config))

		// expose the simple http health endpoint
		if config.httpHealthPath != "" {
			mux.Handle(config.httpHealthPath, httpHealthHandler())
		}
		// expose the hosted FS
		if config.hostedFS != nil {
			var notFoundHandler FSHandler404
			if len(config.hostedFSNotFound) > 0 {
				notFoundHandler = func(w http.ResponseWriter, r *http.Request) (doDefaultFileServe bool) {
					w.WriteHeader(http.StatusNotFound)
					w.Write(config.hostedFSNotFound)
					return
				}
			}
			mux.Handle("/", FileServerWith404(http.FS(config.hostedFS), notFoundHandler))
		}

		httpServer.Handler = InterceptHttp(mux, interceptors...)

		tlsConfig := config.httpTlsConfig
		if tlsConfig == nil {
			tlsConfig = config.grpcTlsConfig
		}
		tlsConfig = tlsConfig.Clone() // so we can modify settings at will
		if tlsConfig != nil {
			tlsConfig.NextProtos = []string{"h2"}
			tlsConfig.ClientAuth = 0 // disable mTLS for web request
		}
		httpServer.TLSConfig = tlsConfig

		httpLis, err := net.Listen("tcp", config.httpAddress)
		if err != nil {
			return fmt.Errorf("http: %w", err)
		}

		addRunMsg(httpLis.Addr(), "HTTP")
		webPageUrl, err = url.Parse("http://" + httpLis.Addr().String())
		if err != nil {
			return fmt.Errorf("parse %v: %w", httpLis.Addr(), err)
		}
		if !config.noGrpcWeb {
			addRunMsg(httpLis.Addr(), "grpc-web")
		}
		go func() { done <- httpServer.Serve(httpLis) }()

		if httpServer.TLSConfig != nil {
			httpsLis, err := net.Listen("tcp", config.httpsAddress)
			if err != nil {
				return fmt.Errorf("https: %w", err)
			}
			addRunMsg(httpsLis.Addr(), "HTTPS")
			webPageUrl, err = url.Parse("https://" + httpsLis.Addr().String())
			if err != nil {
				return fmt.Errorf("parse %v: %w", httpsLis.Addr(), err)
			}
			if !config.noGrpcWeb {
				addRunMsg(httpsLis.Addr(), "grpc-web")
			}
			go func() { done <- httpServer.ServeTLS(httpsLis, "", "") }()
		}
	}

	// print out what we're serving
	var lines []string
	for addr, protocols := range runMsg {
		lines = append(lines, fmt.Sprintf("%v : %v", addr, strings.Join(protocols, ", ")))
	}
	fmt.Printf("Server started serving\n  %v\n", strings.Join(lines, "\n  "))
	if webPageUrl != nil {
		if webPageUrl.Hostname() == "::" {
			if port := webPageUrl.Port(); port != "" {
				webPageUrl.Host = fmt.Sprintf("localhost:%v", webPageUrl.Port())
			} else {
				webPageUrl.Host = "localhost"
			}
		}
		if webPageUrl.Path == "" {
			webPageUrl.Path = "/"
		}
		fmt.Println()
		fmt.Printf("  Admin page:\t\t%v\n", webPageUrl)
		if config.ca != nil {
			webPageUrl.Path = "/__/playground/ca-cert.pem"
			fmt.Printf("  CA Cert:\t\t%v\n", webPageUrl)
			if config.mTLS {
				webPageUrl.Path = "/__/playground/client.pem"
				fmt.Printf("  Client credentials:\t%v  // new per request\n", webPageUrl)
			}
		}
	}

	select {
	case err := <-done:
		// something stopped
		if err != nil {
			return fmt.Errorf("a background task stopped: %w", err)
		} else {
			log.Printf("Shutting down")
		}
	case <-ctx.Done():
		// stop was requested
		log.Printf("Shutting down")
	}

	return nil
}

func loadOrCreateCerts(forceCertGen bool, certCacheDir string) (*ca, *tls.Certificate, error) {
	var cacheDir string
	if !forceCertGen {
		// check for cached certs and use them if we can
		dir := certCacheDir
		if dir == "" {
			dir = filepath.Join(os.TempDir(), CertDir)
		}
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Printf("Unable to load/cache certs: %s\n", err)
		} else {
			cacheDir = dir
		}
	}
	// generate a cert to use
	fmt.Println("Generating self-signed certificates...")
	fmt.Print("  CA...")
	var fromCache bool
	ca, fromCache, err := LoadOrCreateSelfSignedCA(withCacheDir(cacheDir))
	if fromCache {
		fmt.Printf(" loaded from %s!\n", ca.cacheDir)
	} else {
		fmt.Println(" done!")
	}
	if err != nil {
		return nil, nil, fmt.Errorf("ca %w", err)
	}

	fmt.Print("  Server Cert...")
	serverCert, fromCache, err := ca.LoadOrCreateServerCert()
	if fromCache {
		fmt.Printf(" loaded from %s!\n", ca.cacheDir)
	} else {
		fmt.Println(" done!")
	}
	if err != nil {
		return nil, nil, fmt.Errorf("server %w", err)
	}
	return ca, serverCert, err
}

func grpcWebInterceptor(grpcWebServer *grpcweb.WrappedGrpcServer) HttpInterceptor {
	return func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		if grpcWebServer.IsAcceptableGrpcCorsRequest(r) || grpcWebServer.IsGrpcWebRequest(r) {
			next = grpcWebServer
		}

		next.ServeHTTP(w, r)
	}
}

func registerHealth(grpcServer *grpc.Server) {
	healthServer := health.NewServer()
	for name := range grpcServer.GetServiceInfo() {
		healthServer.SetServingStatus(name, grpc_health_v1.HealthCheckResponse_SERVING)
	}
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
}

func httpHealthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "ok, %v %v %v%v %v", r.Method, r.Proto, r.Host, r.RequestURI, r.Header)
	})
}

func configJsonHandler(config *Config) http.Handler {
	configCors := cors.New(cors.Options{
		AllowedMethods: []string{"GET"},
	})
	return configCors.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(config)
		if err != nil {
			log.Printf("ERR: json encode")
		}
	}))
}

func caCertHandler(config *Config) http.Handler {
	if config.ca == nil {
		return http.NotFoundHandler()
	}

	configCors := cors.New(cors.Options{
		AllowedMethods: []string{"GET"},
	})
	return configCors.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/x-x509-ca-cert")
		if err := config.ca.WriteCACertPEM(w); err != nil {
			log.Printf("ca cert encode err after header send: %v", err)
		}
	}))
}

func clientCertHandler(config *Config) http.Handler {
	if config.ca == nil || !config.mTLS {
		return http.NotFoundHandler()
	}

	configCors := cors.New(cors.Options{
		AllowedMethods: []string{"GET"},
	})
	return configCors.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/x-pem")
		log.Printf("Generating new client auth cert...")
		if err := config.ca.WriteClientCert(w); err != nil {
			log.Printf("client cert encode err after header send: %v", err)
		}
	}))
}

type App struct {
	Apis []server.GrpcApi
	Host http.FileSystem

	TlsConfig *tls.Config

	ctx context.Context

	done chan error
	stop func()
}

func NewApp(ctx context.Context) *App {
	done := make(chan error)
	ctx, stop := context.WithCancel(ctx)
	return &App{
		ctx:  ctx,
		done: done,
		stop: stop,
	}
}

func (a *App) WithApis(apis ...server.GrpcApi) {
	a.Apis = append(a.Apis, apis...)
}

func (a *App) HostDir(dir string) {
	a.Host = http.Dir(dir)
}

// WithServerCert causes the gRPC server to only accept tls connections.
// Calling this forces the gRPC server to enable server TLS.
// See: grpc.Creds
func (a *App) WithServerCert(cert tls.Certificate) {
	if a.TlsConfig == nil {
		a.TlsConfig = new(tls.Config)
		a.TlsConfig.NextProtos = []string{"h2"}
	}
	a.TlsConfig.Certificates = append(a.TlsConfig.Certificates, cert)
}

func (a *App) ServeAddress(addr string) error {
	// setup the root tcp listener
	tcpLis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %v %w", addr, err)
	}

	// TCP listening
	tcpMux := cmux.New(tcpLis)
	httpLis := tcpMux.Match(cmux.HTTP1Fast())
	elseLis := tcpMux.Match(cmux.Any())

	// TLS setup
	tlsLis := a.withTls(elseLis)
	tlsMux := cmux.New(tlsLis)
	// grpcTls := tlsLis
	grpcTls := tlsMux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpTls := tlsMux.Match(cmux.Any())

	// grpc setup
	grpcServer := a.SetupGrpcServer(a.Apis...)
	defer grpcServer.GracefulStop()

	// grpc-web setup
	httpServer := SetupHttpServer(grpcServer)
	defer httpServer.Close()

	http.Handle("/health", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer, "ok, %v %v %v%v %v", request.Method, request.Proto, request.Host, request.RequestURI, request.Header)
	}))

	// setup html hosting
	if a.Host != nil {
		http.Handle("/", http.FileServer(a.Host))
	}

	// start everything
	go func() { a.done <- grpcServer.Serve(grpcTls) }()
	go func() { a.done <- httpServer.Serve(httpLis) }()
	go func() { a.done <- httpServer.Serve(httpTls) }()
	go func() { a.done <- tcpMux.Serve() }()
	go func() { a.done <- tlsMux.Serve() }()

	// log what we've done
	var apis []string
	if grpcServer != nil {
		apis = append(apis, "gPRC API")
	}
	if httpServer != nil {
		apis = append(apis, "grpc-web API")
	}
	if a.Host != nil {
		apis = append(apis, "hosted HTML")
	}
	log.Printf("Serving %s on %s", strings.Join(apis, ", "), addr)

	// wait for termination
	select {
	case err := <-a.done:
		// something stopped
		if err != nil {
			return fmt.Errorf("a background task stopped: %w", err)
		} else {
			log.Printf("Shutting down")
		}
	case <-a.ctx.Done():
		// stop was requested
		log.Printf("Shutting down")
	}
	return nil
}

func (a *App) withTls(lis net.Listener) net.Listener {
	if a.TlsConfig == nil {
		return lis
	}

	return tls.NewListener(lis, a.TlsConfig)
}
