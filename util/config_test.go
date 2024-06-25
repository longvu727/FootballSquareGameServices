package util

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (suite *ConfigTestSuite) TestLoadConfig() {
	config, err := LoadConfig(".", "app_test", "env")
	if err != nil {
		suite.Fail("Unable to load config file, " + err.Error())
	}

	suite.Equal("123", config.MySQLDSN)
}

func TestConfigTestSuite(T *testing.T) {
	suite.Run(T, new(ConfigTestSuite))
}
