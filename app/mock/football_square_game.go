// Code generated by MockGen. DO NOT EDIT.
// Source: footballsquaregameservices/app (interfaces: FootballSquareGame)

// Package mockfootballsquaregameapp is a generated GoMock package.
package mockfootballsquaregameapp

import (
	app "footballsquaregameservices/app"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	resources "github.com/longvu727/FootballSquaresLibs/util/resources"
)

// MockFootballSquareGame is a mock of FootballSquareGame interface.
type MockFootballSquareGame struct {
	ctrl     *gomock.Controller
	recorder *MockFootballSquareGameMockRecorder
}

// MockFootballSquareGameMockRecorder is the mock recorder for MockFootballSquareGame.
type MockFootballSquareGameMockRecorder struct {
	mock *MockFootballSquareGame
}

// NewMockFootballSquareGame creates a new mock instance.
func NewMockFootballSquareGame(ctrl *gomock.Controller) *MockFootballSquareGame {
	mock := &MockFootballSquareGame{ctrl: ctrl}
	mock.recorder = &MockFootballSquareGameMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFootballSquareGame) EXPECT() *MockFootballSquareGameMockRecorder {
	return m.recorder
}

// CreateFootballSquareGame mocks base method.
func (m *MockFootballSquareGame) CreateFootballSquareGame(arg0 app.CreateGameParams, arg1 *resources.Resources) (*app.CreateGameResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFootballSquareGame", arg0, arg1)
	ret0, _ := ret[0].(*app.CreateGameResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFootballSquareGame indicates an expected call of CreateFootballSquareGame.
func (mr *MockFootballSquareGameMockRecorder) CreateFootballSquareGame(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFootballSquareGame", reflect.TypeOf((*MockFootballSquareGame)(nil).CreateFootballSquareGame), arg0, arg1)
}

// GetFootballSquareGame mocks base method.
func (m *MockFootballSquareGame) GetFootballSquareGame(arg0 app.GetGameParams, arg1 *resources.Resources) (*app.GetGameResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFootballSquareGame", arg0, arg1)
	ret0, _ := ret[0].(*app.GetGameResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFootballSquareGame indicates an expected call of GetFootballSquareGame.
func (mr *MockFootballSquareGameMockRecorder) GetFootballSquareGame(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFootballSquareGame", reflect.TypeOf((*MockFootballSquareGame)(nil).GetFootballSquareGame), arg0, arg1)
}

// ReserveSquare mocks base method.
func (m *MockFootballSquareGame) ReserveSquare(arg0 app.ReserveSquareParams, arg1 *resources.Resources) (*app.ReserveSquareResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReserveSquare", arg0, arg1)
	ret0, _ := ret[0].(*app.ReserveSquareResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReserveSquare indicates an expected call of ReserveSquare.
func (mr *MockFootballSquareGameMockRecorder) ReserveSquare(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReserveSquare", reflect.TypeOf((*MockFootballSquareGame)(nil).ReserveSquare), arg0, arg1)
}
