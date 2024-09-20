package routes

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type RoutesTestSuite struct {
	suite.Suite
}

func TestRoutesTestSuite(t *testing.T) {
	suite.Run(t, new(RoutesTestSuite))
}

func (suite *RoutesTestSuite) TestRegister() {
	routes := NewRoutes()
	serveMux := routes.Register(nil)

	suite.NotNil(serveMux)
}
