package router

import (
	"compress/flate"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"

	rmGateway "gojo/gateways/rick_and_morty"
	rmHandler "gojo/handlers/rick_and_morty"
)

type ApiRouterConfig struct {
	Handler chi.Router
}

type ApiRouter struct {
	handler chi.Router
}

func NewApiRouter(cfg *ApiRouterConfig) (*ApiRouter, error) {
	switch {
	case cfg == nil:
		return nil, fmt.Errorf("missing config parameter")
	case cfg.Handler == nil:
		return nil, fmt.Errorf("missing Handler parameter")
	}

	return &ApiRouter{
		handler: cfg.Handler,
	}, nil
}

func (r *ApiRouter) Init() {
	port := os.Getenv("PORT")

	r.handler.Use(middleware.Logger)
	r.handler.Use(middleware.Recoverer)
	r.handler.Use(middleware.Compress(flate.DefaultCompression))
	r.handler.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{fmt.Sprintf("http://localhost:%s", port)},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Accept"},
		AllowCredentials: true,
	}))

	r.handler.Use(render.SetContentType(render.ContentTypeJSON))

	rickAndMortyGateway, err := rmGateway.NewGateway(&rmGateway.GatewayConfig{})
	if err != nil {
		log.Fatal(err)
	}

	rickAndMortyHandler, err := rmHandler.NewHandler(&rmHandler.HandlerConfig{
		ApiClient: rickAndMortyGateway,
	})
	if err != nil {
		log.Fatal(err)
	}

	r.handler.Get("/characters/{id}", rickAndMortyHandler.GetCharacter)
	r.handler.Get("/characters/get/{ids}", rickAndMortyHandler.GetCharacters)
	r.handler.Get("/characters/search", rickAndMortyHandler.SearchCharacters)
	r.handler.Get("/characters/list", rickAndMortyHandler.ListCharacters)

	err = http.ListenAndServe(":"+port, r.handler)
	if err != nil {
		log.Fatal(err)
	}
}
