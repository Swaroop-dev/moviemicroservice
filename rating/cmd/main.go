package main

import (
	"log"
	"net/http"

	rating "movieapp.com/rating/internal/controller"
	"movieapp.com/rating/internal/handler/http"
	"movieapp.com/rating/internal/repository/memory"
)

func main() {
	log.Println("starting the rating service")
	repo := memory.New()
	ctrl := rating.New(repo)
	h := httpHandler.New(ctrl)

	http.Handle("/rating", http.HandlerFunc(h.Handle))

	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
