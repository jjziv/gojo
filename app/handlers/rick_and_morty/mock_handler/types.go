// Code generated by MockGen. DO NOT EDIT.
// Source: handlers/rick_and_morty/types.go

// Package mock_rick_and_morty is a generated GoMock package.
package mock_rick_and_morty

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHandler is a mock of Handler interface.
type MockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder
}

// MockHandlerMockRecorder is the mock recorder for MockHandler.
type MockHandlerMockRecorder struct {
	mock *MockHandler
}

// NewMockHandler creates a new mock instance.
func NewMockHandler(ctrl *gomock.Controller) *MockHandler {
	mock := &MockHandler{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandler) EXPECT() *MockHandlerMockRecorder {
	return m.recorder
}

// GetCharacter mocks base method.
func (m *MockHandler) GetCharacter(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetCharacter", w, r)
}

// GetCharacter indicates an expected call of GetCharacter.
func (mr *MockHandlerMockRecorder) GetCharacter(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCharacter", reflect.TypeOf((*MockHandler)(nil).GetCharacter), w, r)
}

// GetCharacters mocks base method.
func (m *MockHandler) GetCharacters(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetCharacters", w, r)
}

// GetCharacters indicates an expected call of GetCharacters.
func (mr *MockHandlerMockRecorder) GetCharacters(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCharacters", reflect.TypeOf((*MockHandler)(nil).GetCharacters), w, r)
}

// ListCharacters mocks base method.
func (m *MockHandler) ListCharacters(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ListCharacters", w, r)
}

// ListCharacters indicates an expected call of ListCharacters.
func (mr *MockHandlerMockRecorder) ListCharacters(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCharacters", reflect.TypeOf((*MockHandler)(nil).ListCharacters), w, r)
}

// SearchCharacters mocks base method.
func (m *MockHandler) SearchCharacters(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SearchCharacters", w, r)
}

// SearchCharacters indicates an expected call of SearchCharacters.
func (mr *MockHandlerMockRecorder) SearchCharacters(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchCharacters", reflect.TypeOf((*MockHandler)(nil).SearchCharacters), w, r)
}
