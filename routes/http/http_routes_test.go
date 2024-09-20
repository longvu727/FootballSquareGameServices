package httproutes

import (
	"bytes"
	"encoding/json"
	"errors"
	"footballsquaregameservices/app"
	mockfootballsquaregameapp "footballsquaregameservices/app/mock"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/longvu727/FootballSquaresLibs/util/resources"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/rand"
)

type HttpRoutesTestSuite struct {
	suite.Suite
}

func TestHttpRoutesTestSuite(t *testing.T) {
	suite.Run(t, new(HttpRoutesTestSuite))
}

func (suite *HttpRoutesTestSuite) getTestError() error {
	return errors.New("test error")
}

func (suite *HttpRoutesTestSuite) TestCreateGame() {

	url := "/CreateGame"

	createGameParams := randomCreateGameParams()
	requestParams, _ := json.Marshal(createGameParams)

	ctrl := gomock.NewController(suite.T())

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestParams))
	suite.NoError(err)

	httpRecorder := httptest.NewRecorder()

	mockFootballSquareGame := mockfootballsquaregameapp.NewMockFootballSquareGame(ctrl)
	mockFootballSquareGame.EXPECT().
		CreateFootballSquareGame(gomock.Any(), gomock.Any()).
		Times(1).
		Return(&app.CreateGameResponse{GameGUID: uuid.NewString()}, nil)

	serveMux := http.NewServeMux()
	httpRoutes := NewHTTPRoutes(mockFootballSquareGame)
	httpRoutes.Register(serveMux, nil)

	handler, pattern := serveMux.Handler(req)
	suite.Equal(http.MethodPost+" "+url, pattern)

	handler.ServeHTTP(httpRecorder, req)

	suite.Equal(httpRecorder.Code, http.StatusOK)
}

func (suite *HttpRoutesTestSuite) TestCreateFootballSquareGameError() {

	url := "/CreateGame"

	createGameParams := randomCreateGameParams()
	requestParams, _ := json.Marshal(createGameParams)

	ctrl := gomock.NewController(suite.T())

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestParams))
	suite.NoError(err)

	httpRecorder := httptest.NewRecorder()

	mockFootballSquareGame := mockfootballsquaregameapp.NewMockFootballSquareGame(ctrl)
	mockFootballSquareGame.EXPECT().
		CreateFootballSquareGame(gomock.Any(), gomock.Any()).
		Times(1).
		Return(&app.CreateGameResponse{}, suite.getTestError())

	serveMux := http.NewServeMux()
	httpRoutes := NewHTTPRoutes(mockFootballSquareGame)
	httpRoutes.Register(serveMux, nil)

	handler, pattern := serveMux.Handler(req)
	suite.Equal(http.MethodPost+" "+url, pattern)

	handler.ServeHTTP(httpRecorder, req)

	suite.Equal(httpRecorder.Code, http.StatusInternalServerError)
}

func (suite *HttpRoutesTestSuite) TestGetGame() {

	game := randomGetGameResponse()
	url := "/GetGame/" + game.GameGUID

	ctrl := gomock.NewController(suite.T())

	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte(``)))
	suite.NoError(err)

	httpRecorder := httptest.NewRecorder()

	mockFootballSquareGame := mockfootballsquaregameapp.NewMockFootballSquareGame(ctrl)
	mockFootballSquareGame.EXPECT().
		GetFootballSquareGame(gomock.Any(), gomock.Any()).
		Times(1).
		Return(game, nil)

	serveMux := http.NewServeMux()
	httpRoutes := NewHTTPRoutes(mockFootballSquareGame)
	httpRoutes.Register(serveMux, nil)

	handler, _ := serveMux.Handler(req)
	handler.ServeHTTP(httpRecorder, req)

	suite.Equal(httpRecorder.Code, http.StatusOK)
}

func (suite *HttpRoutesTestSuite) TestGetGameError() {

	game := randomGetGameResponse()
	url := "/GetGame/" + game.GameGUID

	ctrl := gomock.NewController(suite.T())

	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte(``)))
	suite.NoError(err)

	httpRecorder := httptest.NewRecorder()

	mockFootballSquareGame := mockfootballsquaregameapp.NewMockFootballSquareGame(ctrl)
	mockFootballSquareGame.EXPECT().
		GetFootballSquareGame(gomock.Any(), gomock.Any()).
		Times(1).
		Return(&app.GetGameResponse{}, suite.getTestError())

	serveMux := http.NewServeMux()
	httpRoutes := NewHTTPRoutes(mockFootballSquareGame)
	httpRoutes.Register(serveMux, nil)

	handler, _ := serveMux.Handler(req)
	handler.ServeHTTP(httpRecorder, req)

	suite.Equal(httpRecorder.Code, http.StatusInternalServerError)
}

func (suite *HttpRoutesTestSuite) TestReserveSquare() {

	url := "/ReserveSquares"

	reserveSquareParams := app.ReserveSquareParams{
		GameGUID:    uuid.NewString(),
		UserGUID:    uuid.NewString(),
		RowIndex:    rand.Intn(1000),
		ColumnIndex: rand.Intn(1000),
	}
	requestParams, _ := json.Marshal(reserveSquareParams)

	ctrl := gomock.NewController(suite.T())

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestParams))
	suite.NoError(err)

	httpRecorder := httptest.NewRecorder()

	mockFootballSquareGame := mockfootballsquaregameapp.NewMockFootballSquareGame(ctrl)
	mockFootballSquareGame.EXPECT().
		ReserveSquare(gomock.Any(), gomock.Any()).
		Times(1).
		Return(&app.ReserveSquareResponse{Reserved: true}, nil)

	redisClient, redisMock := redismock.NewClientMock()
	redisMock.ExpectPublish("SquareReserved:"+reserveSquareParams.GameGUID, "SquareReserved").RedisNil()

	resources := resources.Resources{
		RedisClient: redisClient,
	}
	serveMux := http.NewServeMux()
	httpRoutes := NewHTTPRoutes(mockFootballSquareGame)
	httpRoutes.Register(serveMux, &resources)

	handler, pattern := serveMux.Handler(req)
	suite.Equal(http.MethodPost+" "+url, pattern)

	handler.ServeHTTP(httpRecorder, req)

	suite.Equal(httpRecorder.Code, http.StatusOK)
}

func (suite *HttpRoutesTestSuite) TestReserveSquareError() {

	url := "/ReserveSquares"

	reserveSquareParams := app.ReserveSquareParams{
		GameGUID:    uuid.NewString(),
		UserGUID:    uuid.NewString(),
		RowIndex:    rand.Intn(1000),
		ColumnIndex: rand.Intn(1000),
	}
	requestParams, _ := json.Marshal(reserveSquareParams)

	ctrl := gomock.NewController(suite.T())

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestParams))
	suite.NoError(err)

	httpRecorder := httptest.NewRecorder()

	mockFootballSquareGame := mockfootballsquaregameapp.NewMockFootballSquareGame(ctrl)
	mockFootballSquareGame.EXPECT().
		ReserveSquare(gomock.Any(), gomock.Any()).
		Times(1).
		Return(&app.ReserveSquareResponse{}, suite.getTestError())

	serveMux := http.NewServeMux()
	httpRoutes := NewHTTPRoutes(mockFootballSquareGame)
	httpRoutes.Register(serveMux, nil)

	handler, pattern := serveMux.Handler(req)
	suite.Equal(http.MethodPost+" "+url, pattern)

	handler.ServeHTTP(httpRecorder, req)

	suite.Equal(httpRecorder.Code, http.StatusInternalServerError)
}

func (suite *HttpRoutesTestSuite) TestHome() {

	url := "/"
	req, err := http.NewRequest(http.MethodPost, url, nil)
	suite.NoError(err)

	httpRecorder := httptest.NewRecorder()

	ctrl := gomock.NewController(suite.T())
	mockFootballSquareGame := mockfootballsquaregameapp.NewMockFootballSquareGame(ctrl)

	serveMux := http.NewServeMux()
	httpRoutes := NewHTTPRoutes(mockFootballSquareGame)
	httpRoutes.Register(serveMux, nil)

	handler, pattern := serveMux.Handler(req)
	suite.Equal(url, pattern)

	handler.ServeHTTP(httpRecorder, req)

	suite.Equal(httpRecorder.Code, http.StatusOK)
}

func (suite *HttpRoutesTestSuite) TestOptionsMethod() {

	url := "/"
	req, err := http.NewRequest(http.MethodOptions, url, nil)
	suite.NoError(err)

	httpRecorder := httptest.NewRecorder()

	ctrl := gomock.NewController(suite.T())
	mockFootballSquareGame := mockfootballsquaregameapp.NewMockFootballSquareGame(ctrl)

	serveMux := http.NewServeMux()
	httpRoutes := NewHTTPRoutes(mockFootballSquareGame)
	httpRoutes.Register(serveMux, nil)

	handler, pattern := serveMux.Handler(req)
	suite.Equal(url, pattern)

	handler.ServeHTTP(httpRecorder, req)

	suite.Equal(httpRecorder.Code, http.StatusOK)
}

func randomCreateGameParams() app.CreateGameParams {
	return app.CreateGameParams{
		Sport:      "football",
		SquareSize: rand.Int31n(10),
		TeamA:      "red" + strconv.Itoa(rand.Intn(1000)),
		TeamB:      "blue" + strconv.Itoa(rand.Intn(1000)),
	}
}

func randomGetGameResponse() *app.GetGameResponse {
	squareSize := rand.Intn(10)
	rowpoints := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	colpoints := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

	rand.Shuffle(len(rowpoints), func(i, j int) {
		rowpoints[i], rowpoints[j] = rowpoints[j], rowpoints[i]
	})

	rand.Shuffle(len(colpoints), func(i, j int) {
		colpoints[i], colpoints[j] = colpoints[j], colpoints[i]
	})

	return &app.GetGameResponse{
		GameGUID:        uuid.NewString(),
		Sport:           "football",
		TeamA:           "red" + strconv.Itoa(rand.Intn(1000)),
		TeamB:           "red" + strconv.Itoa(rand.Intn(1000)),
		SquareSize:      squareSize,
		RowPoints:       strings.Join(rowpoints, ","),
		ColumnPoints:    strings.Join(colpoints, ","),
		FootballSquares: randomFootballSquares(squareSize),
	}
}

func randomFootballSquares(squareSize int) []app.FootballSquare {
	footballSquares := []app.FootballSquare{}
	for i := 0; i < squareSize*squareSize; i++ {
		winner := false
		winnerQuarterNumber := rand.Intn(3) + 1
		if winnerQuarterNumber > 0 {
			winner = true
		}
		footballSquares = append(footballSquares, app.FootballSquare{
			ColumnIndex:         rand.Intn(10),
			RowIndex:            rand.Intn(10),
			WinnerQuarterNumber: winnerQuarterNumber,
			Winner:              winner,
			UserGUID:            uuid.NewString(),
			UserName:            "user" + strconv.Itoa(rand.Intn(1000)),
			UserAlias:           "u" + strconv.Itoa(rand.Intn(1000)),
		})
	}

	return footballSquares
}
