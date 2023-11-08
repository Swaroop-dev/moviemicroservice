package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// "gopkg.in/yaml.v3"
	"movieapp.com/gen"
	"movieapp.com/metadata/internal/controller/metadata"

	// httphandler "movieapp.com/metadata/internal/handler/http"
	"gopkg.in/yaml.v3"
	grpchandler "movieapp.com/metadata/internal/handler/grpc"
	"movieapp.com/metadata/internal/repository/memory"
	"movieapp.com/pkg/discovery"
	"movieapp.com/pkg/discovery/consul"
)

const serviceName = "metadata"

func main() {
	log.Println("Starting the movie metadata service")
	f, err := os.Open("base.yaml")

	if err != nil {
		panic(err)
	}
	defer f.Close()
	var cfg serviceConfig

	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}
	var port int
	flag.IntVar(&port, "port", 8081, "API handler port")
	flag.Parse()
	log.Printf("Starting the metadata service on port %d", port)
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", cfg.API.Port)); err != nil {
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
	repo := memory.New()
	metadatactrl := metadata.New(repo)
	h := grpchandler.New(metadatactrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterMetadataServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
