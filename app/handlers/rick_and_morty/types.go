package rick_and_morty

import (
	"net/http"

	"gojo/gateways/rick_and_morty"
)

type Handler interface {
	GetCharacter(w http.ResponseWriter, r *http.Request)
	GetCharacters(w http.ResponseWriter, r *http.Request)
	SearchCharacters(w http.ResponseWriter, r *http.Request)
	ListCharacters(w http.ResponseWriter, r *http.Request)
}

type CharacterResponse struct {
	Data rick_and_morty.Character `json:"data,omitempty"`
}

type ListCharactersResponse struct {
	Data []rick_and_morty.Character `json:"data,omitempty"`
}
