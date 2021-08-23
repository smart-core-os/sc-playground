package run

import (
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/smart-core-os/sc-golang/pkg/server"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func (a *App) SetupGrpcServer(apis ...server.GrpcApi) *grpc.Server {
	grpcServer := grpc.NewServer()
	healthApi := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthApi)
	reflection.Register(grpcServer)
	for _, a := range apis {
		a.Register(grpcServer)
	}
	return grpcServer
}

func SetupHttpServer(server *grpc.Server) *http.Server {
	grpcWebServer := grpcweb.WrapServer(server,
		grpcweb.WithOriginFunc(func(origin string) bool {
			return true
		}))
	httpServer := &http.Server{
		Handler: http.HandlerFunc(func(response http.ResponseWriter, req *http.Request) {
			if grpcWebServer.IsAcceptableGrpcCorsRequest(req) || grpcWebServer.IsGrpcWebRequest(req) {
				grpcWebServer.ServeHTTP(response, req)
				return
			}

			// standard http routing
			http.DefaultServeMux.ServeHTTP(response, req)
		}),
	}

	err := http2.ConfigureServer(httpServer, nil)
	if err != nil {
		panic(err)
	}
	return httpServer
}
