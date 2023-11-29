package transport

// import (
// 	"bytes"
// 	"context"
// 	"database/sql"
// 	db "gophermart/internal/database"
// 	"gophermart/internal/mocks"
// 	"gophermart/internal/services"
// 	"gophermart/pkg/logger"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/go-resty/resty/v2"
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/suite"
// )

// type HandlerTestSuite struct {
// 	suite.Suite
// 	client *resty.Client
// 	server *httptest.Server
// }

// func (suite *HandlerTestSuite) SetupSuite() {
// 	suite.client = resty.New()
// }

// func TestSuiteTestJSONHandler(t *testing.T) {
// 	suite.Run(t, &HandlerTestSuite{})
// }

// func (suite *HandlerTestSuite) TearDownSuite() {
// 	suite.server.Close()
// }

// // строка для создания файла с моками
// // mockgen -destination="internal/mocks/mock_store.go" -package=mocks "gophermart/internal/database" StoragerDB
// func (suite *HandlerTestSuite) TestGetUsers() {

// 	ctrl := gomock.NewController(suite.T())
// 	defer ctrl.Finish()

// 	ctx := context.Background()
// 	m := mocks.NewMockStoragerDB(ctrl)
// 	logger, err := logger.NewLogger("Info")
// 	suite.NoError(err)
// 	h := New(ctx, m, logger)
// 	suite.server = httptest.NewServer(http.HandlerFunc(h.Registration))

// 	body := services.UserAuthInfo{Login: "User1", Password: "123"}

// 	// expUser := db.User{uid: "string",
// 	// 	login: "user1",
// 	// 	hash:  "234"}

// 	u := db.User{Uid: "1324", Login: "User1", Hash: "111"}

// 	m.EXPECT().GetUser(ctx, body).Return(u, nil)

// 	url := "/api/user/register"

// 	resp, err := suite.client.R().
// 		SetBody(body).
// 		Post(suite.server.URL + url)

// 	suite.NoError(err)
// 	suite.Equal(200, resp.StatusCode())
// 	// v, err := strconv.Atoi(string(resp.Body()))
// 	suite.NoError(err)
// 	// suite.Equal(, v)

// }

// func TestRegistration(t *testing.T) {
// 	// Test case 1: Successful registration
// 	t.Run("Successful registration", func(t *testing.T) {
// 		// Create a new handlersData instance

// 		ctrl := gomock.NewController(t)
// 		defer ctrl.Finish()

// 		m := mocks.NewMockStoragerDB(ctrl)
// 		logger, _ := logger.NewLogger("Info")
// 		// suite.NoError(err)
// 		h := &handlersData{
// 			storage: m,
// 			logger:  logger,
// 		}

// 		u := db.User{Uid: "1324", Login: "john", Hash: "6b12e30ad2d20c42d3e38e120191224f0852467e6441aa48ef05834d12810c06"}
// 		ctx := context.Background()
// 		body := services.UserAuthInfo{Login: "john", Password: "password123"}

// 		m.EXPECT().GetUser(ctx, body.Login).Return(u, sql.ErrNoRows)
// 		m.EXPECT().AddUser(ctx, u.Login, u.Hash).Return(nil)

// 		// Create a new HTTP request with the registration data
// 		payload := []byte(`{"login": "john", "password": "password123"}`)
// 		req, err := http.NewRequest("POST", "/registration", bytes.NewBuffer(payload))
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		// Create a new HTTP response recorder
// 		rr := httptest.NewRecorder()

// 		// Call the Registration function
// 		h.Registration(rr, req)

// 		// Check the HTTP status code
// 		if rr.Code != http.StatusOK {
// 			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
// 		}

// 		// Check the response headers
// 		expectedHeaders := map[string]string{
// 			"Content-Type": "application/json",
// 		}
// 		for key, value := range expectedHeaders {
// 			if rr.Header().Get(key) != value {
// 				t.Errorf("Expected header %s: %s, got %s", key, value, rr.Header().Get(key))
// 			}
// 		}
// 	})
// }
