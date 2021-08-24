package run

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/rs/cors"
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
	if config.grpcTlsConfig != nil {
		addRunMsg(grpcLis.Addr(), "Secure gRPC")
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(config.grpcTlsConfig)))
	} else if !config.insecure {
		// generate a cert to use
		fmt.Println("Generating self-signed server certificate...")
		fmt.Println("Use --server-certfile and --server-keyfile to provide your own")
		caCert, serverCert, err := genServerCert()
		if err != nil {
			return err
		}
		config.caCert = caCert
		config.grpcTlsConfig = &tls.Config{
			Certificates: []tls.Certificate{
				*serverCert,
			},
		}
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(config.grpcTlsConfig)))
		if err != nil {
			return err
		}
		addRunMsg(grpcLis.Addr(), "Self-signed gRPC")
	} else {
		addRunMsg(grpcLis.Addr(), "Insecure gRPC")
	}
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
		// expose the config used to start the app
		mux.Handle("/__/playground/config.json", configJsonHandler(config))
		mux.Handle("/__/playground/ca-cert.pem", caCertHandler(config))
		// expose the simple http health endpoint
		if config.httpHealthPath != "" {
			mux.Handle(config.httpHealthPath, httpHealthHandler())
		}
		// expose the hosted FS
		if config.hostedFS != nil {
			mux.Handle("/", http.FileServer(http.FS(config.hostedFS)))
		}

		httpServer.Handler = InterceptHttp(mux, interceptors...)

		tlsConfig := config.httpTlsConfig
		if tlsConfig == nil {
			tlsConfig = config.grpcTlsConfig
		}
		tlsConfig = tlsConfig.Clone() // so we can modify settings at will
		if tlsConfig != nil {
			tlsConfig.NextProtos = []string{"h2"}
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
		lines = append(lines, fmt.Sprintf("%v : %v", strings.Join(protocols, ", "), addr))
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
		fmt.Printf("Admin page: %v\n", webPageUrl)
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
	if config.caCert == nil {
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
		err := pem.Encode(w, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: config.caCert.Certificate[0],
		})
		if err != nil {
			log.Printf("ca cert encode err after header send: %v", err)
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
