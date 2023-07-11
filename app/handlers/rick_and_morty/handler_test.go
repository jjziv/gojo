package rick_and_morty

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gojo/gateways/rick_and_morty"
	mockGateway "gojo/gateways/rick_and_morty/mock_gateway"
)

const (
	testCharacterID          = "1"
	testErrorText            = "an error"
	testMultipleCharacterIDs = "1,2"
	testSearchCharacterQuery = "Rick"
)

type errorBody struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func TestHandler_NewHandler(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when no config passed in", func(t *testing.T) {
		t.Parallel()

		_, err := NewHandler(nil)

		assert.EqualError(t, fmt.Errorf("missing config parameter"), err.Error())
	})

	t.Run("it returns an error when no ApiClient passed in", func(t *testing.T) {
		t.Parallel()

		_, err := NewHandler(&HandlerConfig{ApiClient: nil})

		assert.EqualError(t, fmt.Errorf("missing ApiClient parameter"), err.Error())
	})

	t.Run("it successfully returns a Handler", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, err := NewHandler(&HandlerConfig{
			ApiClient: mockGateway.NewMockGateway(ctrl),
		})

		assert.Nil(t, err)
	})
}

func TestHandler_GetCharacter(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when the API Client returns an error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gatewayMock := mockGateway.NewMockGateway(ctrl)

		gatewayMock.EXPECT().GetCharacter(testCharacterID).Return(rick_and_morty.Character{}, fmt.Errorf(testErrorText))

		h, err := NewHandler(&HandlerConfig{
			ApiClient: gatewayMock,
		})

		if err != nil {
			t.FailNow()
		}

		router := chi.NewRouter()
		router.Get("/characters/{id}", h.GetCharacter)

		req, err := http.NewRequest("GET", fmt.Sprintf("/characters/%s", testCharacterID), nil)
		if err != nil {
			t.FailNow()
		}

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		jsonFromRequest, err := io.ReadAll(rec.Body)
		if err != nil {
			t.FailNow()
		}

		response := errorBody{}

		err = json.Unmarshal(jsonFromRequest, &response)
		if err != nil {
			t.FailNow()
		}

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, http.StatusText(http.StatusInternalServerError), response.Status)
		assert.Equal(t, testErrorText, response.Error)
	})

	t.Run("it successfully returns a Character", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedTime, _ := time.Parse(time.RFC3339, "2017-11-04T18:48:46.250Z")

		expectedCharacter := rick_and_morty.Character{
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

		gatewayMock := mockGateway.NewMockGateway(ctrl)

		gatewayMock.EXPECT().GetCharacter(testCharacterID).Return(expectedCharacter, nil)

		h, err := NewHandler(&HandlerConfig{
			ApiClient: gatewayMock,
		})

		if err != nil {
			t.FailNow()
		}

		router := chi.NewRouter()
		router.Get("/characters/{id}", h.GetCharacter)

		req, err := http.NewRequest("GET", fmt.Sprintf("/characters/%s", testCharacterID), nil)
		if err != nil {
			t.FailNow()
		}

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		jsonFromRequest, err := io.ReadAll(rec.Body)
		if err != nil {
			t.FailNow()
		}

		response := CharacterResponse{}

		err = json.Unmarshal(jsonFromRequest, &response)
		if err != nil {
			t.FailNow()
		}

		expectedResponse := CharacterResponse{
			Data: expectedCharacter,
		}

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectedResponse, response)
	})
}

func TestHandler_GetCharacters(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when the API Client returns an error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gatewayMock := mockGateway.NewMockGateway(ctrl)

		gatewayMock.EXPECT().GetCharacters(testMultipleCharacterIDs).Return([]rick_and_morty.Character{}, fmt.Errorf(testErrorText))

		h, err := NewHandler(&HandlerConfig{
			ApiClient: gatewayMock,
		})

		if err != nil {
			t.FailNow()
		}

		router := chi.NewRouter()
		router.Get("/characters/get/{ids}", h.GetCharacters)

		req, err := http.NewRequest("GET", fmt.Sprintf("/characters/get/%s", testMultipleCharacterIDs), nil)
		if err != nil {
			t.FailNow()
		}

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		jsonFromRequest, err := io.ReadAll(rec.Body)
		if err != nil {
			t.FailNow()
		}

		response := errorBody{}

		err = json.Unmarshal(jsonFromRequest, &response)
		if err != nil {
			t.FailNow()
		}

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, http.StatusText(http.StatusInternalServerError), response.Status)
		assert.Equal(t, testErrorText, response.Error)
	})

	t.Run("it successfully returns a list of characters", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedRickTime, _ := time.Parse(time.RFC3339, "2017-11-04T18:48:46.250Z")
		expectedMortyTime, _ := time.Parse(time.RFC3339, "2017-11-04T18:50:21.651Z")

		expectedCharacters := []rick_and_morty.Character{
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

		gatewayMock := mockGateway.NewMockGateway(ctrl)

		gatewayMock.EXPECT().GetCharacters(testMultipleCharacterIDs).Return(expectedCharacters, nil)

		h, err := NewHandler(&HandlerConfig{
			ApiClient: gatewayMock,
		})

		if err != nil {
			t.FailNow()
		}

		router := chi.NewRouter()
		router.Get("/characters/get/{ids}", h.GetCharacters)

		req, err := http.NewRequest("GET", fmt.Sprintf("/characters/get/%s", testMultipleCharacterIDs), nil)
		if err != nil {
			t.FailNow()
		}

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		jsonFromRequest, err := io.ReadAll(rec.Body)
		if err != nil {
			t.FailNow()
		}

		response := ListCharactersResponse{}

		err = json.Unmarshal(jsonFromRequest, &response)
		if err != nil {
			t.FailNow()
		}

		expectedResponse := ListCharactersResponse{
			Data: expectedCharacters,
		}

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectedResponse, response)
	})
}

func TestHandler_SearchCharacters(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when no search parameter is passed in", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h, err := NewHandler(&HandlerConfig{
			ApiClient: mockGateway.NewMockGateway(ctrl),
		})

		if err != nil {
			t.FailNow()
		}

		router := chi.NewRouter()
		router.Get("/characters/search", h.SearchCharacters)

		req, err := http.NewRequest("GET", "/characters/search", nil)
		if err != nil {
			t.FailNow()
		}

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		jsonFromRequest, err := io.ReadAll(rec.Body)
		if err != nil {
			t.FailNow()
		}

		response := errorBody{}

		err = json.Unmarshal(jsonFromRequest, &response)
		if err != nil {
			t.FailNow()
		}

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, http.StatusText(http.StatusBadRequest), response.Status)
		assert.Equal(t, http.StatusText(http.StatusBadRequest), response.Error)
	})

	t.Run("it returns an error when an invalid search parameter is passed in", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h, err := NewHandler(&HandlerConfig{
			ApiClient: mockGateway.NewMockGateway(ctrl),
		})

		if err != nil {
			t.FailNow()
		}

		router := chi.NewRouter()
		router.Get("/characters/search", h.SearchCharacters)

		req, err := http.NewRequest("GET", "/characters/search?name=001", nil)
		if err != nil {
			t.FailNow()
		}

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		jsonFromRequest, err := io.ReadAll(rec.Body)
		if err != nil {
			t.FailNow()
		}

		response := errorBody{}

		err = json.Unmarshal(jsonFromRequest, &response)
		if err != nil {
			t.FailNow()
		}

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, http.StatusText(http.StatusBadRequest), response.Status)
		assert.Equal(t, http.StatusText(http.StatusBadRequest), response.Error)
	})

	t.Run("it returns an error when the API Client returns an error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gatewayMock := mockGateway.NewMockGateway(ctrl)

		gatewayMock.EXPECT().SearchCharacters(testSearchCharacterQuery).Return([]rick_and_morty.Character{}, fmt.Errorf(testErrorText))

		h, err := NewHandler(&HandlerConfig{
			ApiClient: gatewayMock,
		})

		if err != nil {
			t.FailNow()
		}

		router := chi.NewRouter()
		router.Get("/characters/search", h.SearchCharacters)

		req, err := http.NewRequest("GET", fmt.Sprintf("/characters/search?name=%s", testSearchCharacterQuery), nil)
		if err != nil {
			t.FailNow()
		}

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		jsonFromRequest, err := io.ReadAll(rec.Body)
		if err != nil {
			t.FailNow()
		}

		response := errorBody{}

		err = json.Unmarshal(jsonFromRequest, &response)
		if err != nil {
			t.FailNow()
		}

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, http.StatusText(http.StatusInternalServerError), response.Status)
		assert.Equal(t, testErrorText, response.Error)
	})

	t.Run("it successfully returns a list of characters", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedTime, _ := time.Parse(time.RFC3339, "2017-11-04T18:48:46.250Z")
		expectedCharacter := rick_and_morty.Character{
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

		gatewayMock := mockGateway.NewMockGateway(ctrl)

		gatewayMock.EXPECT().SearchCharacters(testSearchCharacterQuery).Return([]rick_and_morty.Character{expectedCharacter}, nil)

		h, err := NewHandler(&HandlerConfig{
			ApiClient: gatewayMock,
		})

		if err != nil {
			t.FailNow()
		}

		router := chi.NewRouter()
		router.Get("/characters/search", h.SearchCharacters)

		req, err := http.NewRequest("GET", fmt.Sprintf("/characters/search?name=%s", testSearchCharacterQuery), nil)
		if err != nil {
			t.FailNow()
		}

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		jsonFromRequest, err := io.ReadAll(rec.Body)
		if err != nil {
			t.FailNow()
		}

		response := ListCharactersResponse{}

		err = json.Unmarshal(jsonFromRequest, &response)
		if err != nil {
			t.FailNow()
		}

		expectedResponse := ListCharactersResponse{
			Data: []rick_and_morty.Character{expectedCharacter},
		}

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectedResponse, response)
	})
}

func TestHandler_ListCharacters(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when the API Client returns an error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gatewayMock := mockGateway.NewMockGateway(ctrl)

		gatewayMock.EXPECT().ListCharacters().Return([]rick_and_morty.Character{}, fmt.Errorf(testErrorText))

		h, err := NewHandler(&HandlerConfig{
			ApiClient: gatewayMock,
		})

		if err != nil {
			t.FailNow()
		}

		router := chi.NewRouter()
		router.Get("/characters/list", h.ListCharacters)

		req, err := http.NewRequest("GET", "/characters/list", nil)
		if err != nil {
			t.FailNow()
		}

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		jsonFromRequest, err := io.ReadAll(rec.Body)
		if err != nil {
			t.FailNow()
		}

		response := errorBody{}

		err = json.Unmarshal(jsonFromRequest, &response)
		if err != nil {
			t.FailNow()
		}

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, http.StatusText(http.StatusInternalServerError), response.Status)
		assert.Equal(t, testErrorText, response.Error)
	})

	t.Run("it successfully returns a list of characters", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedRickTime, _ := time.Parse(time.RFC3339, "2017-11-04T18:48:46.250Z")
		expectedMortyTime, _ := time.Parse(time.RFC3339, "2017-11-04T18:50:21.651Z")

		expectedCharacters := []rick_and_morty.Character{
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

		gatewayMock := mockGateway.NewMockGateway(ctrl)

		gatewayMock.EXPECT().ListCharacters().Return(expectedCharacters, nil)

		h, err := NewHandler(&HandlerConfig{
			ApiClient: gatewayMock,
		})

		if err != nil {
			t.FailNow()
		}

		router := chi.NewRouter()
		router.Get("/characters/list", h.ListCharacters)

		req, err := http.NewRequest("GET", "/characters/list", nil)
		if err != nil {
			t.FailNow()
		}

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		jsonFromRequest, err := io.ReadAll(rec.Body)
		if err != nil {
			t.FailNow()
		}

		response := ListCharactersResponse{}

		err = json.Unmarshal(jsonFromRequest, &response)
		if err != nil {
			t.FailNow()
		}

		expectedResponse := ListCharactersResponse{
			Data: expectedCharacters,
		}

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectedResponse, response)
	})
}
