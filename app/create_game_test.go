package app

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mockdb "github.com/longvu727/FootballSquaresLibs/DB/db/mock"
	"github.com/longvu727/FootballSquaresLibs/services"
	mockservices "github.com/longvu727/FootballSquaresLibs/services/mock"
	"github.com/longvu727/FootballSquaresLibs/util"
	"github.com/longvu727/FootballSquaresLibs/util/resources"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/rand"
)

type CreateGameTestSuite struct {
	suite.Suite
}

func TestCreateGameTestSuite(t *testing.T) {
	suite.Run(t, new(CreateGameTestSuite))
}

func (suite *CreateGameTestSuite) TestCreateGame() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	game := randomGame()

	mockServices := mockservices.NewMockServices(ctrl)

	mockServices.EXPECT().
		CreateSquare(gomock.Any(), gomock.Any()).
		Times(1).
		Return(randomSquare(), nil)
	
	mockServices.EXPECT().
		CreateGame(gomock.Any(), gomock.Any()).
		Times(1).
		Return(game, nil)
	
	mockServices.EXPECT().
		CreateFootballSquareGame(gomock.Any(), gomock.Any()).
		Times(1).
		Return(randomFootballSquareGame(), nil)

	config, err := util.LoadConfig("../env", "app", "json")
	suite.NoError(err)

	resources := resources.NewResources(config, mockdb.NewMockMySQL(ctrl), mockServices, context.Background())

	createSquareParams := CreateGameParams{
		Sport:      "football",
		SquareSize: rand.Int31n(10),
		TeamA:      "red",
		TeamB:      "blue",
	}

	createFootballSquareGameResponse, err := NewFootballSquareGameApp().CreateFootballSquareGame(createSquareParams, resources)
	suite.NoError(err)

	suite.Equal(createFootballSquareGameResponse.GameGUID, game.GameGUID)
	suite.Greater(len(createFootballSquareGameResponse.ToJson()), 0)
}

func (suite *CreateGameTestSuite) TestCreateGameCreateFootballSquareGameError() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	game := randomGame()

	mockServices := mockservices.NewMockServices(ctrl)

	mockServices.EXPECT().
		CreateSquare(gomock.Any(), gomock.Any()).
		Times(1).
		Return(randomSquare(), nil)
	
	mockServices.EXPECT().
		CreateGame(gomock.Any(), gomock.Any()).
		Times(1).
		Return(game, nil)
	
	mockServices.EXPECT().
		CreateFootballSquareGame(gomock.Any(), gomock.Any()).
		Times(1).
		Return(services.CreateFootballSquareGameResponse{}, errors.New("test error"))

	config, err := util.LoadConfig("../env", "app", "json")
	suite.NoError(err)

	resources := resources.NewResources(config, mockdb.NewMockMySQL(ctrl), mockServices, context.Background())

	_, err = NewFootballSquareGameApp().CreateFootballSquareGame(CreateGameParams{}, resources)
	suite.Error(err)
}

func (suite *CreateGameTestSuite) TestCreateGameCreateGameError() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	mockServices := mockservices.NewMockServices(ctrl)

	mockServices.EXPECT().
		CreateSquare(gomock.Any(), gomock.Any()).
		Times(1).
		Return(randomSquare(), nil)
	
	mockServices.EXPECT().
		CreateGame(gomock.Any(), gomock.Any()).
		Times(1).
		Return(services.CreateGameResponse{}, errors.New("test error"))

	config, err := util.LoadConfig("../env", "app", "json")
	suite.NoError(err)

	resources := resources.NewResources(config, mockdb.NewMockMySQL(ctrl), mockServices, context.Background())

	_, err = NewFootballSquareGameApp().CreateFootballSquareGame(CreateGameParams{}, resources)
	suite.Error(err)
}

func (suite *CreateGameTestSuite) TestCreateGameCreateSquareError() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	mockServices := mockservices.NewMockServices(ctrl)

	mockServices.EXPECT().
		CreateSquare(gomock.Any(), gomock.Any()).
		Times(1).
		Return(services.CreateSquareResponse{}, errors.New("test error"))

	config, err := util.LoadConfig("../env", "app", "json")
	suite.NoError(err)

	resources := resources.NewResources(config, mockdb.NewMockMySQL(ctrl), mockServices, context.Background())

	_, err = NewFootballSquareGameApp().CreateFootballSquareGame(CreateGameParams{}, resources)
	suite.Error(err)
}

func randomGame() services.CreateGameResponse {
	return services.CreateGameResponse{GameGUID: uuid.NewString(), GameID: rand.Int63n(1000)}
}

func randomSquare() services.CreateSquareResponse {
	return services.CreateSquareResponse{SquareID: rand.Intn(1000)}
}

func randomFootballSquareGame() services.CreateFootballSquareGameResponse {
	footballSquaresGameIDs := []int64{}

	squareSize := randomSquareSize()
	for i := 1; i <= squareSize*squareSize; i++ {
		footballSquaresGameIDs = append(footballSquaresGameIDs, rand.Int63n(1000))
	}

	return services.CreateFootballSquareGameResponse{
		FootballSquaresGameIDs: footballSquaresGameIDs,
	}
}

func randomSquareSize() int {
	return rand.Intn(9) + 1
}
