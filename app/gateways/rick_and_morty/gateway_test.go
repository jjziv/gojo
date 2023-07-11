package rick_and_morty

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const (
	testCharacterID          = "1"
	testErrorText            = "an error"
	testMultipleCharacterIDs = "1,2"
	testSearchCharacterQuery = "Rick"
)

func TestGateway_NewGateway(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when no config passed in", func(t *testing.T) {
		t.Parallel()

		_, err := NewGateway(nil)

		assert.EqualError(t, fmt.Errorf("missing config parameter"), err.Error())
	})

	t.Run("it successfully returns a Gateway", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, err := NewGateway(&GatewayConfig{})

		assert.Nil(t, err)
	})
}

func TestGateway_GetCharacter(t *testing.T) {
	t.Run("it returns an error if the API returns an error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g, err := NewGateway(&GatewayConfig{})
		if err != nil {
			t.FailNow()
		}

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", baseURI+"character/"+testCharacterID,
			func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf(testErrorText)
			})

		result, err := g.GetCharacter(testCharacterID)

		assert.Equal(t, Character{}, result)
		assert.Equal(t, "Get \"https://rickandmortyapi.com/api/character/1\": an error", err.Error())
	})

	t.Run("it returns an error if decoding the API response returns an error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g, err := NewGateway(&GatewayConfig{})
		if err != nil {
			t.FailNow()
		}

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", baseURI+"character/"+testCharacterID,
			func(req *http.Request) (*http.Response, error) {
				resp, err := httpmock.NewJsonResponse(200, `{"Name": "Foo`+"\u001a"+`"}`)
				if err != nil {
					return httpmock.NewStringResponse(500, ""), nil
				}
				return resp, nil
			})

		result, err := g.GetCharacter(testCharacterID)

		assert.Equal(t, Character{}, result)
		assert.Error(t, err)
	})

	t.Run("it successfully returns a Character", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g, err := NewGateway(&GatewayConfig{})
		if err != nil {
			t.FailNow()
		}

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		expectedTime, _ := time.Parse(time.RFC3339, "2017-11-04T18:48:46.250Z")
		expectedCharacter := Character{
			Id:      1,
			Name:    "Rick Sanchez",
			Status:  "Alive",
			Species: "Human",
			Type:    "",
			Gender:  "Male",
			Image:   "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
			Url:     "https://rickandmortyapi.com/api/character/1",
			Created: expectedTime,
		}

		httpmock.RegisterResponder("GET", baseURI+"character/"+testCharacterID,
			func(req *http.Request) (*http.Response, error) {
				resp, err := httpmock.NewJsonResponse(200, expectedCharacter)
				if err != nil {
					return httpmock.NewStringResponse(500, ""), nil
				}
				return resp, nil
			})

		result, err := g.GetCharacter(testCharacterID)

		assert.Equal(t, expectedCharacter, result)
		assert.Nil(t, err)
	})
}

func TestGateway_GetCharacters(t *testing.T) {
	t.Run("it returns an error if the API returns an error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g, err := NewGateway(&GatewayConfig{})
		if err != nil {
			t.FailNow()
		}

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", baseURI+"character/"+testMultipleCharacterIDs,
			func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf(testErrorText)
			})

		result, err := g.GetCharacters(testMultipleCharacterIDs)

		assert.Equal(t, []Character{}, result)
		assert.Equal(t, "Get \"https://rickandmortyapi.com/api/character/1,2\": an error", err.Error())
	})

	t.Run("it returns an error if decoding the API response returns an error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g, err := NewGateway(&GatewayConfig{})
		if err != nil {
			t.FailNow()
		}

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", baseURI+"character/"+testMultipleCharacterIDs,
			func(req *http.Request) (*http.Response, error) {
				resp, err := httpmock.NewJsonResponse(200, `{"Name": "Foo`+"\u001a"+`"}`)
				if err != nil {
					return httpmock.NewStringResponse(500, ""), nil
				}
				return resp, nil
			})

		result, err := g.GetCharacters(testMultipleCharacterIDs)

		assert.Equal(t, []Character{}, result)
		assert.Error(t, err)
	})

	t.Run("it successfully returns a list of characters", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g, err := NewGateway(&GatewayConfig{})
		if err != nil {
			t.FailNow()
		}

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		expectedRickTime, _ := time.Parse(time.RFC3339, "2017-11-04T18:48:46.250Z")
		expectedMortyTime, _ := time.Parse(time.RFC3339, "2017-11-04T18:50:21.651Z")

		expectedCharacters := []Character{
			{
				Id:      1,
				Name:    "Rick Sanchez",
				Status:  "Alive",
				Species: "Human",
				Type:    "",
				Gender:  "Male",
				Image:   "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
				Url:     "https://rickandmortyapi.com/api/character/1",
				Created: expectedRickTime,
			},
			{
				Id:      2,
				Name:    "Morty Smith",
				Status:  "Alive",
				Species: "Human",
				Type:    "",
				Gender:  "Male",
				Image:   "https://rickandmortyapi.com/api/character/avatar/2.jpeg",
				Url:     "https://rickandmortyapi.com/api/character/2",
				Created: expectedMortyTime,
			},
		}

		httpmock.RegisterResponder("GET", baseURI+"character/"+testMultipleCharacterIDs,
			func(req *http.Request) (*http.Response, error) {
				resp, err := httpmock.NewJsonResponse(200, expectedCharacters)
				if err != nil {
					return httpmock.NewStringResponse(500, ""), nil
				}
				return resp, nil
			})

		result, err := g.GetCharacters(testMultipleCharacterIDs)

		assert.Equal(t, expectedCharacters, result)
		assert.Nil(t, err)
	})
}

func TestGateway_SearchCharacters(t *testing.T) {
	t.Run("it returns an error if the API returns an error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g, err := NewGateway(&GatewayConfig{})
		if err != nil {
			t.FailNow()
		}

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", baseURI+"character?name="+testSearchCharacterQuery,
			func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf(testErrorText)
			})

		result, err := g.SearchCharacters(testSearchCharacterQuery)

		assert.Equal(t, []Character{}, result)
		assert.Equal(t, "Get \"https://rickandmortyapi.com/api/character?name=Rick\": an error", err.Error())
	})

	t.Run("it returns an error if decoding the API response returns an error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g, err := NewGateway(&GatewayConfig{})
		if err != nil {
			t.FailNow()
		}

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", baseURI+"character?name="+testSearchCharacterQuery,
			func(req *http.Request) (*http.Response, error) {
				resp, err := httpmock.NewJsonResponse(200, `{"Name": "Foo`+"\u001a"+`"}`)
				if err != nil {
					return httpmock.NewStringResponse(500, ""), nil
				}
				return resp, nil
			})

		result, err := g.SearchCharacters(testSearchCharacterQuery)

		assert.Equal(t, []Character{}, result)
		assert.Error(t, err)
	})

	t.Run("it successfully returns a list of characters", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g, err := NewGateway(&GatewayConfig{})
		if err != nil {
			t.FailNow()
		}

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		expectedTime, _ := time.Parse(time.RFC3339, "2017-11-04T18:48:46.250Z")
		expectedCharacter := Character{
			Id:      1,
			Name:    "Rick Sanchez",
			Status:  "Alive",
			Species: "Human",
			Type:    "",
			Gender:  "Male",
			Image:   "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
			Url:     "https://rickandmortyapi.com/api/character/1",
			Created: expectedTime,
		}

		expectedResponse := CharactersListResponse{
			Info: ApiInfo{
				Count: 1,
				Pages: 1,
			},
			Results: []Character{expectedCharacter},
		}

		httpmock.RegisterResponder("GET", baseURI+"character?name="+testSearchCharacterQuery,
			func(req *http.Request) (*http.Response, error) {
				resp, err := httpmock.NewJsonResponse(200, expectedResponse)
				if err != nil {
					return httpmock.NewStringResponse(500, ""), nil
				}
				return resp, nil
			})

		result, err := g.SearchCharacters(testSearchCharacterQuery)

		assert.Equal(t, []Character{expectedCharacter}, result)
		assert.Nil(t, err)
	})
}

func TestGateway_ListCharacters(t *testing.T) {
	t.Run("it returns an error if the API returns an error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g, err := NewGateway(&GatewayConfig{})
		if err != nil {
			t.FailNow()
		}

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", baseURI+"character",
			func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf(testErrorText)
			})

		result, err := g.ListCharacters()

		assert.Equal(t, []Character{}, result)
		assert.Equal(t, "Get \"https://rickandmortyapi.com/api/character\": an error", err.Error())
	})

	t.Run("it returns an error if decoding the API response returns an error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g, err := NewGateway(&GatewayConfig{})
		if err != nil {
			t.FailNow()
		}

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", baseURI+"character",
			func(req *http.Request) (*http.Response, error) {
				resp, err := httpmock.NewJsonResponse(200, `{"Name": "Foo`+"\u001a"+`"}`)
				if err != nil {
					return httpmock.NewStringResponse(500, ""), nil
				}
				return resp, nil
			})

		result, err := g.ListCharacters()

		assert.Equal(t, []Character{}, result)
		assert.Error(t, err)
	})

	t.Run("it successfully returns a list of characters", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g, err := NewGateway(&GatewayConfig{})
		if err != nil {
			t.FailNow()
		}

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		expectedTime, _ := time.Parse(time.RFC3339, "2017-11-04T18:48:46.250Z")
		expectedCharacter := Character{
			Id:      1,
			Name:    "Rick Sanchez",
			Status:  "Alive",
			Species: "Human",
			Type:    "",
			Gender:  "Male",
			Image:   "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
			Url:     "https://rickandmortyapi.com/api/character/1",
			Created: expectedTime,
		}

		expectedResponse := CharactersListResponse{
			Info: ApiInfo{
				Count: 1,
				Pages: 1,
			},
			Results: []Character{expectedCharacter},
		}

		httpmock.RegisterResponder("GET", baseURI+"character",
			func(req *http.Request) (*http.Response, error) {
				resp, err := httpmock.NewJsonResponse(200, expectedResponse)
				if err != nil {
					return httpmock.NewStringResponse(500, ""), nil
				}
				return resp, nil
			})

		result, err := g.ListCharacters()

		assert.Equal(t, []Character{expectedCharacter}, result)
		assert.Nil(t, err)
	})
}
