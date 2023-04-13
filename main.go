package main

import (
	"context"
	"log"

	"github.com/dimat/volume-airports-finder/finder"
	"github.com/dimat/volume-airports-finder/server"
	"github.com/dimat/volume-airports-finder/service"
	"github.com/dimat/volume-airports-finder/utils"
)

func main() {
	srv := server.New(":8080")

	srv.Register("/calculate", service.HandlerFunc[service.FindPathRequest, service.FindPathResponse](
		&service.Finder{PathFinder: finder.NewPathFinder()}))

	ctx := utils.ContextWithSignal(context.Background())
	err := srv.Start(ctx)
	if err != nil {
		log.Fatal("Error running server:", err)
	}
}
