package rick_and_morty

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"gojo/gateways/rick_and_morty"
	"gojo/utilities"
)

type HandlerConfig struct {
	ApiClient rick_and_morty.Gateway
}

type handler struct {
	apiClient rick_and_morty.Gateway
}

func NewHandler(cfg *HandlerConfig) (Handler, error) {
	switch {
	case cfg == nil:
		return nil, fmt.Errorf("missing config parameter")
	case cfg.ApiClient == nil:
		return nil, fmt.Errorf("missing ApiClient parameter")
	}

	return &handler{
		apiClient: cfg.ApiClient,
	}, nil
}

func (h *handler) GetCharacter(w http.ResponseWriter, r *http.Request) {
	characterID := chi.URLParam(r, "id")

	if characterID == "" {
		log.Println("no characterID parameter passed in!")
		utilities.RenderHTTPError(w, r)
		return
	}

	character, err := h.apiClient.GetCharacter(characterID)
	if err != nil {
		log.Println(err)
		utilities.RenderServerError(w, r, err)
		return
	}

	response := CharacterResponse{
		Data: character,
	}

	render.JSON(w, r, response)
}

func (h *handler) GetCharacters(w http.ResponseWriter, r *http.Request) {
	characterIDs := chi.URLParam(r, "ids")

	characterList, err := h.apiClient.GetCharacters(characterIDs)
	if err != nil {
		log.Println(err)
		utilities.RenderServerError(w, r, err)
		return
	}

	response := ListCharactersResponse{
		Data: characterList,
	}

	render.JSON(w, r, response)
}

func (h *handler) SearchCharacters(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	if name == "" {
		log.Println("no name parameter passed in!")
		utilities.RenderHTTPError(w, r)
		return
	}

	reg := regexp.MustCompile(`[^A-Za-z]`)
	searchParameter := reg.ReplaceAllString(name, "")
	if searchParameter == "" {
		log.Println("invalid name parameter passed in!")
		utilities.RenderHTTPError(w, r)
		return
	}

	characterList, err := h.apiClient.SearchCharacters(searchParameter)
	if err != nil {
		log.Println(err)
		utilities.RenderServerError(w, r, err)
		return
	}

	response := ListCharactersResponse{
		Data: characterList,
	}

	render.JSON(w, r, response)
}

func (h *handler) ListCharacters(w http.ResponseWriter, r *http.Request) {
	characterList, err := h.apiClient.ListCharacters()
	if err != nil {
		log.Println(err)
		utilities.RenderServerError(w, r, err)
		return
	}

	response := ListCharactersResponse{
		Data: characterList,
	}

	render.JSON(w, r, response)
}
