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
	"gopkg.in/yaml.v3"
	"movieapp.com/pkg/discovery"
	"movieapp.com/pkg/discovery/consul"
	rating "movieapp.com/rating/internal/controller/rating"

	//httpHandler "movieapp.com/rating/internal/handler/http"
	"movieapp.com/gen"
	grpcHandler "movieapp.com/rating/internal/handler/grpc"
	"movieapp.com/rating/internal/repository/mysql"
)

const serviceName = "rating"

func main() {
	log.Println("Starting rating service")
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
	flag.IntVar(&port, "port", 8082, "API handler port")
	flag.Parse()
	log.Printf("Starting the rating service on port %d", port)
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
	repo, err := mysql.New()
	if err != nil {
		panic(err)
	}
	ctrl := rating.New(repo, nil)
	h := grpcHandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))

	// http.Handle("/rating", http.HandlerFunc(h.Handle))

	// if err := http.ListenAndServe(":8082", nil); err != nil {
	// 	panic(err)
	// }

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterRatingServiceServer(srv, h)

	if err := srv.Serve(lis); err != nil {
		panic(err)
	}

}
