package main

import (
	"log"
	"net/http"

	"movieapp.com/metadata/internal/controller/metadata"
	httphandler "movieapp.com/metadata/internal/handler/http"
	"movieapp.com/metadata/internal/repository/memory"
)

func main() {
	log.Println("starting the movie service")
	repo := memory.New()
	metadarepo := metadata.New(repo)
	h := httphandler.New(metadarepo)

	http.Handle("/metadata", http.HandlerFunc(h.GetMetaData))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
