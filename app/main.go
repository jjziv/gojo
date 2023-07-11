package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	applicationPort = "1337"
)

func main() {
	apiRouter, err := NewApiRouter(&ApiRouterConfig{
		Handler: chi.NewRouter(),
	})
	if err != nil {
		log.Fatal(err)
	}

	apiRouter.Init()

	err = http.ListenAndServe(":"+applicationPort, apiRouter.handler)
	if err != nil {
		log.Fatal(err)
	}
}
