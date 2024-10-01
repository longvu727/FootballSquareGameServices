package websocketroutes

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"footballsquaregameservices/app"
	mockfootballsquaregameapp "footballsquaregameservices/app/mock"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/longvu727/FootballSquaresLibs/util/resources"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/rand"
)

type WebsocketRoutesTestSuite struct {
	suite.Suite
}

func TestWebsocketRoutesTestSuite(t *testing.T) {
	suite.Run(t, new(WebsocketRoutesTestSuite))
}

func (suite *WebsocketRoutesTestSuite) getTestError() error {
	return errors.New("test error")
}

func (suite *WebsocketRoutesTestSuite) TestSubscribeGameUpgraderError() {

	game := randomGetGameResponse()
	url := "/Subscribe/GetGame/" + game.GameGUID

	ctrl := gomock.NewController(suite.T())

	mockFootballSquareGame := mockfootballsquaregameapp.NewMockFootballSquareGame(ctrl)

	resources := resources.Resources{
		Context: context.Background(),
	}

	serveMux := http.NewServeMux()
	wsRoutes := NewWebSocketRoutes(mockFootballSquareGame)
	wsRoutes.Register(serveMux, &resources)

	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte{}))
	suite.NoError(err)

	httpRecorder := httptest.NewRecorder()

	serveMux.ServeHTTP(httpRecorder, req)

	_, _, err = websocket.DefaultDialer.Dial("ws://"+url, nil)
	fmt.Println(err)
	suite.Error(err)
}

func (suite *WebsocketRoutesTestSuite) TestSubscribeGameGetFootballSquareGameError() {

	game := randomGetGameResponse()
	url := "/Subscribe/GetGame/" + game.GameGUID

	ctrl := gomock.NewController(suite.T())

	mockFootballSquareGame := mockfootballsquaregameapp.NewMockFootballSquareGame(ctrl)
	mockFootballSquareGame.EXPECT().
		GetFootballSquareGame(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(&app.GetGameResponse{}, suite.getTestError())

	mRedis, err := miniredis.Run()
	suite.NoError(err)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     mRedis.Addr(),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	resources := resources.Resources{
		RedisClient: redisClient,
		Context:     context.Background(),
	}

	serveMux := http.NewServeMux()
	wsRoutes := NewWebSocketRoutes(mockFootballSquareGame)
	wsRoutes.Register(serveMux, &resources)

	server := httptest.NewServer(serveMux)
	defer server.Close()

	wsUrl := "ws" + strings.TrimPrefix(server.URL, "http") + url

	client, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	suite.NoError(err)
	defer client.Close()

	_, response, _ := client.ReadMessage()
	suite.Equal("", string(response))
}

func (suite *WebsocketRoutesTestSuite) TestSubscribeGameSquareReservedReceiveMessageError() {

	game := randomGetGameResponse()
	url := "/Subscribe/GetGame/" + game.GameGUID

	ctrl := gomock.NewController(suite.T())

	mockFootballSquareGame := mockfootballsquaregameapp.NewMockFootballSquareGame(ctrl)
	mockFootballSquareGame.EXPECT().
		GetFootballSquareGame(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(game, nil)

	mRedis := miniredis.RunT(suite.T())

	redisClient := redis.NewClient(&redis.Options{
		Addr:     mRedis.Addr(),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	resources := resources.Resources{
		RedisClient: redisClient,
		Context:     context.Background(),
	}

	serveMux := http.NewServeMux()
	wsRoutes := NewWebSocketRoutes(mockFootballSquareGame)
	wsRoutes.Register(serveMux, &resources)

	server := httptest.NewServer(serveMux)
	defer server.Close()

	wsUrl := "ws" + strings.TrimPrefix(server.URL, "http") + url

	client, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	suite.NoError(err)
	defer client.Close()

	var gameResponse *app.GetGameResponse

	_, response, _ := client.ReadMessage()
	json.NewDecoder(bytes.NewReader(response)).Decode(&gameResponse)
	suite.Equal(game, gameResponse)

	err = redisClient.Publish(resources.Context, "SquareReserved:"+game.GameGUID+"NotEqual", "SquareReserved").Err()
	suite.NoError(err)

}

func (suite *WebsocketRoutesTestSuite) TestSubscribeGameReloadFootballSquareGameError() {

	game := randomGetGameResponse()
	url := "/Subscribe/GetGame/" + game.GameGUID

	ctrl := gomock.NewController(suite.T())

	mockFootballSquareGame := mockfootballsquaregameapp.NewMockFootballSquareGame(ctrl)
	mockFootballSquareGame.EXPECT().
		GetFootballSquareGame(gomock.Any(), gomock.Any()).
		Times(1).
		Return(game, nil)

	mockFootballSquareGame.EXPECT().
		GetFootballSquareGame(gomock.Any(), gomock.Any()).
		Times(1).
		Return(&app.GetGameResponse{}, suite.getTestError())

	mRedis, err := miniredis.Run()
	suite.NoError(err)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     mRedis.Addr(),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	resources := resources.Resources{
		RedisClient: redisClient,
		Context:     context.Background(),
	}

	serveMux := http.NewServeMux()
	wsRoutes := NewWebSocketRoutes(mockFootballSquareGame)
	wsRoutes.Register(serveMux, &resources)

	server := httptest.NewServer(serveMux)
	defer server.Close()

	wsUrl := "ws" + strings.TrimPrefix(server.URL, "http") + url

	client, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	suite.NoError(err)
	defer client.Close()

	var gameResponse *app.GetGameResponse

	_, response, _ := client.ReadMessage()
	json.NewDecoder(bytes.NewReader(response)).Decode(&gameResponse)
	suite.Equal(game, gameResponse)

	subscribeGameChannel := redisClient.Subscribe(resources.Context, "SubscribeGame:"+game.GameGUID)
	defer subscribeGameChannel.Close()

	err = redisClient.Publish(resources.Context, "SquareReserved:"+game.GameGUID, "SquareReserved").Err()
	suite.NoError(err)
}

func (suite *WebsocketRoutesTestSuite) TestSubscribeGame() {

	game := randomGetGameResponse()
	url := "/Subscribe/GetGame/" + game.GameGUID

	ctrl := gomock.NewController(suite.T())

	mockFootballSquareGame := mockfootballsquaregameapp.NewMockFootballSquareGame(ctrl)
	mockFootballSquareGame.EXPECT().
		GetFootballSquareGame(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(game, nil)

	mRedis, err := miniredis.Run()
	suite.NoError(err)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     mRedis.Addr(),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	resources := resources.Resources{
		RedisClient: redisClient,
		Context:     context.Background(),
	}

	serveMux := http.NewServeMux()
	wsRoutes := NewWebSocketRoutes(mockFootballSquareGame)
	wsRoutes.Register(serveMux, &resources)

	server := httptest.NewServer(serveMux)
	defer server.Close()

	wsUrl := "ws" + strings.TrimPrefix(server.URL, "http") + url

	client, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	suite.NoError(err)
	defer client.Close()

	var gameResponse *app.GetGameResponse

	_, response, _ := client.ReadMessage()
	json.NewDecoder(bytes.NewReader(response)).Decode(&gameResponse)
	suite.Equal(game, gameResponse)

	subscribeGameChannel := redisClient.Subscribe(resources.Context, "SubscribeGame:"+game.GameGUID)
	defer subscribeGameChannel.Close()

	err = redisClient.Publish(resources.Context, "SquareReserved:"+game.GameGUID, "SquareReserved").Err()
	suite.NoError(err)

	subscribeGameresponse, err := subscribeGameChannel.ReceiveMessage(resources.Context)
	suite.NoError(err)
	suite.Equal(string(game.ToJson()), subscribeGameresponse.Payload)
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
