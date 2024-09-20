package app

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/longvu727/FootballSquaresLibs/DB/db/mock"
	"github.com/longvu727/FootballSquaresLibs/services"
	mockservices "github.com/longvu727/FootballSquaresLibs/services/mock"
	"github.com/longvu727/FootballSquaresLibs/util"
	"github.com/longvu727/FootballSquaresLibs/util/resources"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/rand"
)

type ReserveSquareTestSuite struct {
	suite.Suite
}

func TestReserveSquareTestSuite(t *testing.T) {
	suite.Run(t, new(ReserveSquareTestSuite))
}

func (suite *ReserveSquareTestSuite) getTestError() error {
	return errors.New("test error")
}

func (suite *ReserveSquareTestSuite) TestReserveSquare() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	game := randomGameByGUID()
	user := randomUserByGUID()

	mockServices := mockservices.NewMockServices(ctrl)

	mockServices.EXPECT().
		GetGameByGUID(gomock.Any(), gomock.Any()).
		Times(1).
		Return(game, nil)

	mockServices.EXPECT().
		GetUserByGUID(gomock.Any(), gomock.Any()).
		Times(1).
		Return(user, nil)

	mockServices.EXPECT().
		ReserveFootballSquare(gomock.Any(), gomock.Any()).
		Times(1).
		Return(services.ReserveFootballSquareResponse{Reserved: true}, nil)

	config, err := util.LoadConfig("../env", "app", "json")
	suite.NoError(err)

	resources := resources.NewResources(config, mockdb.NewMockMySQL(ctrl), mockServices, context.Background())

	reserveSquareParams := ReserveSquareParams{
		GameGUID:    game.GameGUID,
		UserGUID:    user.UserGUID,
		RowIndex:    rand.Intn(1000),
		ColumnIndex: rand.Intn(1000),
	}

	reserveFootballSquareGameResponse, err := NewFootballSquareGameApp().ReserveSquare(reserveSquareParams, resources)
	suite.NoError(err)

	suite.True(reserveFootballSquareGameResponse.Reserved)
	suite.Greater(len(reserveFootballSquareGameResponse.ToJson()), 0)
}

func (suite *ReserveSquareTestSuite) TestReserveSquareReserveFootballSquareError() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	game := randomGameByGUID()
	user := randomUserByGUID()

	mockServices := mockservices.NewMockServices(ctrl)

	mockServices.EXPECT().
		GetGameByGUID(gomock.Any(), gomock.Any()).
		Times(1).
		Return(game, nil)

	mockServices.EXPECT().
		GetUserByGUID(gomock.Any(), gomock.Any()).
		Times(1).
		Return(user, nil)

	mockServices.EXPECT().
		ReserveFootballSquare(gomock.Any(), gomock.Any()).
		Times(1).
		Return(services.ReserveFootballSquareResponse{}, suite.getTestError())

	config, err := util.LoadConfig("../env", "app", "json")
	suite.NoError(err)

	resources := resources.NewResources(config, mockdb.NewMockMySQL(ctrl), mockServices, context.Background())

	reserveSquareParams := ReserveSquareParams{
		GameGUID:    game.GameGUID,
		UserGUID:    user.UserGUID,
		RowIndex:    rand.Intn(1000),
		ColumnIndex: rand.Intn(1000),
	}

	_, err = NewFootballSquareGameApp().ReserveSquare(reserveSquareParams, resources)
	suite.Error(err)
}

func (suite *ReserveSquareTestSuite) TestReserveSquareGetUserByGUIDError() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	game := randomGameByGUID()
	user := randomUserByGUID()

	mockServices := mockservices.NewMockServices(ctrl)

	mockServices.EXPECT().
		GetGameByGUID(gomock.Any(), gomock.Any()).
		Times(1).
		Return(game, nil)

	mockServices.EXPECT().
		GetUserByGUID(gomock.Any(), gomock.Any()).
		Times(1).
		Return(services.GetUserByGUIDResponse{}, suite.getTestError())

	config, err := util.LoadConfig("../env", "app", "json")
	suite.NoError(err)

	resources := resources.NewResources(config, mockdb.NewMockMySQL(ctrl), mockServices, context.Background())

	reserveSquareParams := ReserveSquareParams{
		GameGUID:    game.GameGUID,
		UserGUID:    user.UserGUID,
		RowIndex:    rand.Intn(1000),
		ColumnIndex: rand.Intn(1000),
	}

	_, err = NewFootballSquareGameApp().ReserveSquare(reserveSquareParams, resources)
	suite.Error(err)
}


func (suite *ReserveSquareTestSuite) TestReserveSquareGetGameByGUIDError() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	game := randomGameByGUID()
	user := randomUserByGUID()

	mockServices := mockservices.NewMockServices(ctrl)

	mockServices.EXPECT().
		GetGameByGUID(gomock.Any(), gomock.Any()).
		Times(1).
		Return(services.GetGameByGUIDResponse{}, suite.getTestError())

	config, err := util.LoadConfig("../env", "app", "json")
	suite.NoError(err)

	resources := resources.NewResources(config, mockdb.NewMockMySQL(ctrl), mockServices, context.Background())

	reserveSquareParams := ReserveSquareParams{
		GameGUID:    game.GameGUID,
		UserGUID:    user.UserGUID,
		RowIndex:    rand.Intn(1000),
		ColumnIndex: rand.Intn(1000),
	}

	_, err = NewFootballSquareGameApp().ReserveSquare(reserveSquareParams, resources)
	suite.Error(err)
}

func randomUserByGUID() services.GetUserByGUIDResponse {
	return services.GetUserByGUIDResponse{
		User: randomUser(),
	}
}
