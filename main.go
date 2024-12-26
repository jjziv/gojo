package main

import (
	"github.com/go-chi/chi/v5"
	"log"

	"gojo/router"
)

func main() {
	apiRouter, err := router.NewApiRouter(&router.ApiRouterConfig{
		Handler: chi.NewRouter(),
	})
	if err != nil {
		log.Fatal(err)
	}

	apiRouter.Init()
}
