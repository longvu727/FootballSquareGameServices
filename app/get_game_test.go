package app

import (
	"context"
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/longvu727/FootballSquaresLibs/DB/db"
	mockdb "github.com/longvu727/FootballSquaresLibs/DB/db/mock"
	"github.com/longvu727/FootballSquaresLibs/util"
	"github.com/longvu727/FootballSquaresLibs/util/resources"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/rand"
)

type GetGameTestSuite struct {
	suite.Suite
}

func (suite *GetGameTestSuite) SetupTest() {
}

func (suite *GetGameTestSuite) TestGetGameByGUID() {
	randomGame := randomGameByGUID()

	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	mockMySQL := mockdb.NewMockMySQL(ctrl)

	mockMySQL.EXPECT().
		GetGameByGUID(gomock.Any(), gomock.Eq(randomGame.GameGuid)).
		Times(1).
		Return(randomGame, nil)

	config, err := util.LoadConfig("../env", "app", "env")
	suite.NoError(err)

	resources := resources.NewResources(config, mockMySQL, context.Background())

	getGameParams := GetGameParams{GameGUID: randomGame.GameGuid}
	game, err := NewGameApp().GetFootballSquareGame(getGameParams, resources)
	suite.NoError(err)

	suite.Equal(randomGame.GameGuid, game.GameGUID)
	suite.Equal(randomGame.Sport.String, game.Sport)
	suite.Equal(randomGame.TeamA.String, game.TeamA)
	suite.Equal(randomGame.TeamB.String, game.TeamB)
}
func randomGame() db.GetGameRow {
	return db.GetGameRow{
		GameID:   rand.Int31n(1000),
		GameGuid: uuid.NewString(),
		Sport:    sql.NullString{String: "football", Valid: true},
		TeamA:    sql.NullString{String: "TeamA", Valid: true},
		TeamB:    sql.NullString{String: "TeamB", Valid: true},
	}
}

func randomGameByGUID() db.GetGameByGUIDRow {
	return db.GetGameByGUIDRow{
		GameID:   rand.Int31n(1000),
		GameGuid: uuid.NewString(),
		Sport:    sql.NullString{String: "football", Valid: true},
		TeamA:    sql.NullString{String: "TeamA", Valid: true},
		TeamB:    sql.NullString{String: "TeamB", Valid: true},
	}
}

func TestGetGameTestSuite(t *testing.T) {
	suite.Run(t, new(GetGameTestSuite))
}
