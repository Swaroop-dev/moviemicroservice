package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"movieapp.com/movie/internal/controller/movie"
	metadata_gateway "movieapp.com/movie/internal/gateway/metadata/http"
	rating_gateway "movieapp.com/movie/internal/gateway/rating/http"
	httpHandler "movieapp.com/movie/internal/handler/http"

	"movieapp.com/pkg/discovery"
	"movieapp.com/pkg/discovery/consul"
)

const serviceName = "movie"

func main() {
	log.Println("starting the movie service ")

	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()
	log.Printf("Starting the movie service on port %d", port)
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	metadatagateway := metadata_gateway.New(registry)
	ratinggateway := rating_gateway.New(registry)
	ctrl := movie.New(ratinggateway, metadatagateway)

	h := httpHandler.New(ctrl)

	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}
