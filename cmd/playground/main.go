package main

import (
	"context"
	"flag"
	"log"

	"github.com/smart-core-os/sc-playground/pkg/apis"
	"github.com/smart-core-os/sc-playground/pkg/run"
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
		apis.PowerSupplyApi(),
	)
	return app.ServeAddress(*bind)
}
