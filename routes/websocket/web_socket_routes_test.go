package websocketroutes

import (
	"bytes"
	"context"
	"fmt"
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
	"github.com/gorilla/websocket"
	"github.com/longvu727/FootballSquaresLibs/util/resources"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/rand"
)

type WebsocketRoutesTestSuite struct {
	suite.Suite
}

func TestWebsocketRoutesTestSuite(t *testing.T) {
	suite.Run(t, new(WebsocketRoutesTestSuite))
}

/*
func (suite *WebsocketRoutesTestSuite) getTestError() error {
	return errors.New("test error")
}
*/

func (suite *WebsocketRoutesTestSuite) TestSubscribeGame() {

	game := randomGetGameResponse()
	url := "/Subscribe/GetGame/" + game.GameGUID

	ctrl := gomock.NewController(suite.T())

	httpRecorder := httptest.NewRecorder()

	mockFootballSquareGame := mockfootballsquaregameapp.NewMockFootballSquareGame(ctrl)
	mockFootballSquareGame.EXPECT().
		GetFootballSquareGame(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(game, nil)

	redisClient, redisServerMock := redismock.NewClientMock()
	//redisServerMock.ExpectPublish("SquareReserved:"+game.GameGUID, `SquareReserved`)
	redisServerMock.ExpectPublish("SubscribeGame:"+game.GameGUID, game.ToJson())

	resources := resources.Resources{
		RedisClient: redisClient,
		Context:     context.Background(),
	}

	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte(``)))
	suite.NoError(err)

	serveMux := http.NewServeMux()
	wsRoutes := NewWebSocketRoutes(mockFootballSquareGame)
	wsRoutes.Register(serveMux, &resources)
	handler, _ := serveMux.Handler(req)

	server := httptest.NewServer(handler)
	defer server.Close()

	wsUrl := "ws" + strings.TrimPrefix(server.URL, "http") + url

	client, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		suite.Fail("could not connect to WebSocket server: %v", err)
	}
	defer client.Close()

	_, response, err := client.ReadMessage()
	if err != nil {
		suite.Fail("could not read message from WebSocket server: %v", err)
	}
	fmt.Println(string(response))

	suite.Equal(httpRecorder.Code, http.StatusOK)
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
