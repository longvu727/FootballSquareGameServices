package app

import (
	"context"
	"errors"
	"strconv"
	"strings"
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

type GetGameTestSuite struct {
	suite.Suite
}

func TestGetGameTestSuite(t *testing.T) {
	suite.Run(t, new(GetGameTestSuite))
}

func (suite *GetGameTestSuite) getTestError() error {
	return errors.New("test error")
}

func (suite *GetGameTestSuite) TestGetGame() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	game := randomGameByGUID()

	mockServices := mockservices.NewMockServices(ctrl)

	mockServices.EXPECT().
		GetGameByGUID(gomock.Any(), gomock.Any()).
		Times(1).
		Return(game, nil)

	mockServices.EXPECT().
		GetFootballSquareGameByGameID(gomock.Any(), gomock.Any()).
		Times(1).
		Return(randomFootballSquareGameByGameID(), nil)

	mockServices.EXPECT().
		GetSquare(gomock.Any(), gomock.Any()).
		Times(1).
		Return(randomGetSquare(), nil)

	mockServices.EXPECT().
		GetUser(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(randomGetUser(), nil)

	config, err := util.LoadConfig("../env", "app", "json")
	suite.NoError(err)

	resources := resources.NewResources(config, mockdb.NewMockMySQL(ctrl), mockServices, context.Background())

	getSquareParams := GetGameParams{
		GameGUID: game.GameGUID,
	}

	getFootballSquareGameResponse, err := NewFootballSquareGameApp().GetFootballSquareGame(getSquareParams, resources)
	suite.NoError(err)

	suite.Equal(getFootballSquareGameResponse.GameGUID, game.GameGUID)
	suite.Greater(len(getFootballSquareGameResponse.ToJson()), 0)
}

func (suite *GetGameTestSuite) TestGetGameGetUserError() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	game := randomGameByGUID()

	mockServices := mockservices.NewMockServices(ctrl)

	mockServices.EXPECT().
		GetGameByGUID(gomock.Any(), gomock.Any()).
		Times(1).
		Return(game, nil)

	mockServices.EXPECT().
		GetFootballSquareGameByGameID(gomock.Any(), gomock.Any()).
		Times(1).
		Return(randomFootballSquareGameByGameID(), nil)

	mockServices.EXPECT().
		GetSquare(gomock.Any(), gomock.Any()).
		Times(1).
		Return(randomGetSquare(), nil)

	mockServices.EXPECT().
		GetUser(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(services.GetUserResponse{}, suite.getTestError())

	config, err := util.LoadConfig("../env", "app", "json")
	suite.NoError(err)

	resources := resources.NewResources(config, mockdb.NewMockMySQL(ctrl), mockServices, context.Background())

	getSquareParams := GetGameParams{
		GameGUID: game.GameGUID,
	}

	getFootballSquareGameResponse, err := NewFootballSquareGameApp().GetFootballSquareGame(getSquareParams, resources)
	suite.NoError(err)

	suite.Equal(getFootballSquareGameResponse.GameGUID, game.GameGUID)
	suite.Greater(len(getFootballSquareGameResponse.ToJson()), 0)
}

func (suite *GetGameTestSuite) TestGetGameGetSquareError() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	game := randomGameByGUID()

	mockServices := mockservices.NewMockServices(ctrl)

	mockServices.EXPECT().
		GetGameByGUID(gomock.Any(), gomock.Any()).
		Times(1).
		Return(game, nil)

	mockServices.EXPECT().
		GetFootballSquareGameByGameID(gomock.Any(), gomock.Any()).
		Times(1).
		Return(randomFootballSquareGameByGameID(), nil)

	mockServices.EXPECT().
		GetSquare(gomock.Any(), gomock.Any()).
		Times(1).
		Return(services.GetSquareResponse{}, suite.getTestError())

	mockServices.EXPECT().
		GetUser(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(services.GetUserResponse{}, suite.getTestError())

	config, err := util.LoadConfig("../env", "app", "json")
	suite.NoError(err)

	resources := resources.NewResources(config, mockdb.NewMockMySQL(ctrl), mockServices, context.Background())

	getSquareParams := GetGameParams{
		GameGUID: game.GameGUID,
	}

	_, err = NewFootballSquareGameApp().GetFootballSquareGame(getSquareParams, resources)
	suite.Error(err)

}

func (suite *GetGameTestSuite) TestGetGameGetFootballSquareGameError() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	game := randomGameByGUID()

	mockServices := mockservices.NewMockServices(ctrl)

	mockServices.EXPECT().
		GetGameByGUID(gomock.Any(), gomock.Any()).
		Times(1).
		Return(game, nil)

	mockServices.EXPECT().
		GetFootballSquareGameByGameID(gomock.Any(), gomock.Any()).
		Times(1).
		Return(services.GetFootballSquareGameByGameIDResponse{}, suite.getTestError())

	config, err := util.LoadConfig("../env", "app", "json")
	suite.NoError(err)

	resources := resources.NewResources(config, mockdb.NewMockMySQL(ctrl), mockServices, context.Background())

	getSquareParams := GetGameParams{
		GameGUID: game.GameGUID,
	}

	_, err = NewFootballSquareGameApp().GetFootballSquareGame(getSquareParams, resources)
	suite.Error(err)
}

func (suite *GetGameTestSuite) TestGetGameGetGameByGUIDError() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	game := randomGameByGUID()

	mockServices := mockservices.NewMockServices(ctrl)

	mockServices.EXPECT().
		GetGameByGUID(gomock.Any(), gomock.Any()).
		Times(1).
		Return(services.GetGameByGUIDResponse{}, suite.getTestError())

	config, err := util.LoadConfig("../env", "app", "json")
	suite.NoError(err)

	resources := resources.NewResources(config, mockdb.NewMockMySQL(ctrl), mockServices, context.Background())

	getSquareParams := GetGameParams{
		GameGUID: game.GameGUID,
	}

	_, err = NewFootballSquareGameApp().GetFootballSquareGame(getSquareParams, resources)
	suite.Error(err)
}

func randomFootballSquareGameElement() services.FootballSquareGameElement {
	return services.FootballSquareGameElement{
		GameID:               rand.Intn(1000),
		FootballSquareGameID: rand.Intn(1000),
		SquareID:             rand.Intn(1000),
		UserID:               rand.Intn(1000),
		Winner:               false,
		WinnerQuarterNumber:  rand.Intn(1),
		RowIndex:             rand.Intn(10),
		ColumnIndex:          rand.Intn(10),
		UserName:             "user" + strconv.Itoa(rand.Intn(1000)),
		UserAlias:            "u" + strconv.Itoa(rand.Intn(1000)),
	}
}

func randomFootballSquareGameByGameID() services.GetFootballSquareGameByGameIDResponse {
	getFootballSquareGameElement := []services.FootballSquareGameElement{}

	squareSize := rand.Intn(99) + 1
	for i := 0; i < squareSize; i++ {
		randomFootballSquareGame := randomFootballSquareGameElement()
		getFootballSquareGameElement = append(getFootballSquareGameElement, randomFootballSquareGame)
	}

	return services.GetFootballSquareGameByGameIDResponse{
		FootballSquares: getFootballSquareGameElement,
	}
}

func randomUser() services.User{
	return services.User{
		UserID:     rand.Intn(1000),
		UserGUID:   uuid.NewString(),
		IP:         strconv.Itoa(rand.Intn(255)) + "." + strconv.Itoa(rand.Intn(255)) + "." + strconv.Itoa(rand.Intn(255)),
		DeviceName: "random device " + strconv.Itoa(rand.Intn(1000)),
		UserName:   "user" + strconv.Itoa(rand.Intn(1000)),
		Alias:      "u" + strconv.Itoa(rand.Intn(1000)),
	}
}

func randomGetUser() services.GetUserResponse {
	return services.GetUserResponse{
		User: randomUser(),
	}
}

func randomGameByGUID() services.GetGameByGUIDResponse {
	return services.GetGameByGUIDResponse{
		Game: services.Game{
			Sport: "football",
			TeamA: "red",
			TeamB: "blue",
		},
	}
}

func randomGetSquare() services.GetSquareResponse {
	rowpoints := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	colpoints := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

	rand.Shuffle(len(rowpoints), func(i, j int) {
		rowpoints[i], rowpoints[j] = rowpoints[j], rowpoints[i]
	})

	rand.Shuffle(len(colpoints), func(i, j int) {
		colpoints[i], colpoints[j] = colpoints[j], colpoints[i]
	})

	return services.GetSquareResponse{
		Square: services.Square{
			SquareGUID:   uuid.NewString(),
			SquareID:     rand.Intn(1000),
			SquareSize:   rand.Intn(10),
			RowPoints:    strings.Join(rowpoints, ","),
			ColumnPoints: strings.Join(rowpoints, ","),
		},
	}
}
