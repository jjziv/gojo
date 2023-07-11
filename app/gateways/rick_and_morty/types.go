package rick_and_morty

import "time"

type Gateway interface {
	GetCharacter(id string) (Character, error)
	GetCharacters(ids string) ([]Character, error)
	SearchCharacters(name string) ([]Character, error)
	ListCharacters() ([]Character, error)
}

type Character struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Species string `json:"species"`
	Type    string `json:"type"`
	Gender  string `json:"gender"`
	Origin  struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"origin"`
	Location struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"location"`
	Image   string    `json:"image"`
	Episode []string  `json:"episode"`
	Url     string    `json:"url"`
	Created time.Time `json:"created"`
}

type ApiInfo struct {
	Count int     `json:"count,omitempty"`
	Pages int     `json:"pages,omitempty"`
	Next  string  `json:"next,omitempty"`
	Prev  *string `json:"prev,omitempty"`
}

type CharactersListResponse struct {
	Info    ApiInfo
	Results []Character `json:"results,omitempty"`
}
