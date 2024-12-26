package rick_and_morty

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURI = "https://rickandmortyapi.com/api/"

type GatewayConfig struct {
	//HttpClient utilities.HttpClient // If using a custom HTTP Client.
}

type gateway struct {
	//httpClient utilities.HttpClient
}

func NewGateway(cfg *GatewayConfig) (Gateway, error) {
	switch {
	case cfg == nil:
		return nil, fmt.Errorf("missing config parameter")
		//case cfg.HttpClient == nil:
		//	return nil, fmt.Errorf("missing HttpClient parameter")
	}

	return &gateway{
		//httpClient: cfg.HttpClient,
	}, nil
}

func (g *gateway) GetCharacter(id string) (Character, error) {
	apiResponse, err := http.Get(baseURI + "character/" + id)
	if err != nil {
		return Character{}, err
	}

	apiData := Character{}

	err = json.NewDecoder(apiResponse.Body).Decode(&apiData)
	if err != nil {
		return Character{}, err
	}

	return apiData, nil
}

func (g *gateway) GetCharacters(ids string) ([]Character, error) {
	apiResponse, err := http.Get(baseURI + "character/" + ids)
	if err != nil {
		return []Character{}, err
	}

	var apiData []Character

	err = json.NewDecoder(apiResponse.Body).Decode(&apiData)
	if err != nil {
		return []Character{}, err
	}

	return apiData, nil
}

func (g *gateway) SearchCharacters(name string) ([]Character, error) {
	characterList, err := g.getAllData(baseURI + "character?name=" + name)
	if err != nil {
		return []Character{}, err
	}

	return characterList, nil
}

func (g *gateway) ListCharacters() ([]Character, error) {
	characterList, err := g.getAllData(baseURI + "character")
	if err != nil {
		return []Character{}, err
	}

	return characterList, nil
}

func (g *gateway) getAllData(url string) ([]Character, error) {
	var allData []Character

	apiResponse, err := http.Get(url)
	if err != nil {
		return allData, err
	}

	apiData := CharactersListResponse{}

	err = json.NewDecoder(apiResponse.Body).Decode(&apiData)
	if err != nil {
		return []Character{}, err
	}

	totalPages := apiData.Info.Pages
	if totalPages == 1 {
		return apiData.Results, nil
	}

	if apiData.Results != nil {
		for _, c := range apiData.Results {
			allData = append(allData, c)
		}
	}

	nextBatch := ""
	if apiData.Info.Next != "" {
		nextBatch = apiData.Info.Next
	}

	for i := 1; i < totalPages; i++ {
		apiResponse, err = http.Get(nextBatch)
		if err != nil {
			return allData, err
		}

		apiData = CharactersListResponse{}

		err = json.NewDecoder(apiResponse.Body).Decode(&apiData)
		if err != nil {
			continue
		}

		if apiData.Results != nil {
			for _, c := range apiData.Results {
				allData = append(allData, c)
			}
		}

		nextBatch = apiData.Info.Next
	}

	return allData, nil
}
