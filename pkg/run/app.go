package run

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/smart-core-os/sc-golang/pkg/server"
	"github.com/soheilhy/cmux"
)

type App struct {
	Apis []server.GrpcApi
	Host http.FileSystem

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

func (a *App) ServeAddress(addr string) error {
	// network setup
	sharedLis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %v %w", addr, err)
	}
	mux := cmux.New(sharedLis)

	// grpc setup
	grpcServer := SetupGrpcServer(a.Apis...)
	defer grpcServer.GracefulStop()

	// grpc-web setup
	httpServer := SetupHttpServer(grpcServer)
	defer httpServer.Close()

	// start everything
	grpcLis := mux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	go func() { a.done <- grpcServer.Serve(grpcLis) }()

	httpLis := mux.Match(cmux.Any())
	go func() { a.done <- httpServer.Serve(httpLis) }()

	// setup html hosting
	if a.Host != nil {
		http.Handle("/", http.FileServer(a.Host))
	}
	go func() { a.done <- mux.Serve() }()

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
