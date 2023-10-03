package main

import (
	"log"
	"net/http"

	"movieapp.com/movie/internal/controller/movie"
	metadata_gateway "movieapp.com/movie/internal/gateway/metadata/http"
	rating_gateway "movieapp.com/movie/internal/gateway/rating/http"
	httpHandler "movieapp.com/movie/internal/handler/http"
)

func main() {
	log.Println("starting the movie service ")

	metadatagateway := metadata_gateway.New("localhost:8081")
	ratinggateway := rating_gateway.New("localhost:8082")
	ctrl := movie.New(ratinggateway, metadatagateway)

	h := httpHandler.New(ctrl)

	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}
