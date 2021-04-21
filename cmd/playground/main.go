package main

import (
	"context"
	"flag"
	"log"

	"git.vanti.co.uk/smartcore/sc-playground/pkg/apis"
	"git.vanti.co.uk/smartcore/sc-playground/pkg/run"
)

var (
	bind = flag.String("bind", ":23557", "grpc server bind")
)

func main() {
	if err := Run(); err != nil {
		log.Printf("Exiting: %v", err)
	}
}

func Run() error {
	flag.Parse()
	app := run.NewApp(context.Background())
	app.WithApis(
		apis.BookingApi(),
		apis.OccupancyApi(),
		apis.OnOffApi(),
	)
	return app.ServeAddress(*bind)
}
