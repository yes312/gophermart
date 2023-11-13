package transport

import (
	"gophermart/internal/services"
	"net/http/httptest"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite
	client *resty.Client
	server *httptest.Server
}

func (suite *HandlerTestSuite) SetupSuite() {
	suite.client = resty.New()
	// h := New(storage)
	// suite.server = httptest.NewServer(http.HandlerFunc(h.Registration))
}

func (suite *HandlerTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *HandlerTestSuite) TestGetUsers() {

	url := "/api/user/register"
	body := services.UserAuthInfo{Login: "User1", Password: "123"}

	resp, err := suite.client.R().
		SetBody(body).
		Post(suite.server.URL + url)

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode())

}
