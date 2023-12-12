package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type tSuite struct {
	suite.Suite
	storage *Storage
}

func TestSuiteTest(t *testing.T) {
	suite.Run(t, &tSuite{})
}

func (ts *tSuite) TestAddUser() {

	ts.T().Log("Тест AddUser()")
	ctx := context.Background()
	ts.TruncateAllTables(ctx)

	type userTestCase struct {
		name        string
		user        User
		expectedErr bool
	}

	tests := []userTestCase{
		{
			name:        "создание пользователя",
			user:        User{Login: "Jhon", Hash: "123"},
			expectedErr: false,
		},
		{
			name:        "повторное создание пользователя",
			user:        User{Login: "Jhon", Hash: "123"},
			expectedErr: true,
		},
		{
			name:        "пустой хеш",
			user:        User{Login: "Jhon", Hash: ""},
			expectedErr: true,
		},
		{
			name:        "пустой логин",
			user:        User{Login: "", Hash: "123"},
			expectedErr: true,
		},
	}

	for _, tc := range tests {
		var isErr bool
		_, err := ts.storage.WithRetry(ctx, ts.storage.AddUser(ctx, tc.user.Login, tc.user.Hash))
		if err != nil {
			isErr = true
		}
		ts.Equal(tc.expectedErr, isErr, tc.name)
	}

}

func (ts *tSuite) TestGetUser() {

	ts.T().Log("Тест GetUser()")
	ctx := context.Background()
	ts.TruncateAllTables(ctx)

	expectedUser := User{
		Login: "Jhon",
		Hash:  "123",
	}
	_, err := ts.storage.WithRetry(ctx, ts.storage.GetUser(ctx, expectedUser.Login))
	ts.True(errors.Is(err, sql.ErrNoRows), "пользователь не существует")

	_, err = ts.storage.WithRetry(ctx, ts.storage.AddUser(ctx, expectedUser.Login, expectedUser.Hash))
	ts.NoError(err)

	userInterface, err := ts.storage.WithRetry(ctx, ts.storage.GetUser(ctx, expectedUser.Login))
	ts.NoError(err)

	user, ok := userInterface.(User)
	ts.True(ok)
	ts.Equal(expectedUser, user)

}

func (ts *tSuite) TestAddOrder() {

	ts.T().Log("Тест TestAddOrder()")
	ctx := context.Background()
	ts.TruncateAllTables(ctx)

	expectedUser := User{
		Login: "Jhon",
		Hash:  "123",
	}

	_, err := ts.storage.WithRetry(ctx, ts.storage.AddUser(ctx, expectedUser.Login, expectedUser.Hash))
	ts.NoError(err)

	type testCase struct {
		name        string
		userID      string
		orderNumber string
		expectedErr bool
		expected    OrderUserID
	}

	testCases := []testCase{
		{
			name:        "добавляем новый заказ",
			orderNumber: "123",
			userID:      "Jhon",
			expectedErr: false,
			expected:    OrderUserID{},
		},
		{
			name:        "добавляем уже существующий заказ",
			orderNumber: "123",
			userID:      "Jhon",
			expectedErr: false,
			expected: OrderUserID{
				OrderNumber: "123",
				UserID:      "Jhon",
			},
		},
		{
			name:        "пустой пользователь",
			orderNumber: "234",
			userID:      "",
			expectedErr: true,
		},
		{
			name:        "пустой номер заказа",
			orderNumber: "",
			userID:      "Jhon",
			expectedErr: true,
		},
		{
			name:        "такой номер заказа есть у другого пользователя",
			orderNumber: "123",
			userID:      "Pharhad",
			expectedErr: false,
			expected: OrderUserID{
				OrderNumber: "123",
				UserID:      "Jhon",
			},
		},
	}

	for _, tc := range testCases {
		var isErr bool
		orderUserIDInterface, err := ts.storage.WithRetry(ctx, ts.storage.AddOrder(ctx, tc.orderNumber, tc.userID))
		orderUserID, _ := orderUserIDInterface.(OrderUserID)
		log.Println("orderUserID", orderUserID)
		if err != nil {
			isErr = true
		}
		ts.Equal(tc.expectedErr, isErr, tc.name)
		ts.Equal(tc.expected, orderUserID, tc.name)
	}

}

func (ts *tSuite) TestPutStatuses() {

	ts.T().Log("Тест TestPutStatuses()")
	ctx := context.Background()
	ts.TruncateAllTables(ctx)

	// подготавливаем данные для теста
	expectedUser := User{Login: "Jhon", Hash: "123"}
	_, err := ts.storage.WithRetry(ctx, ts.storage.AddUser(ctx, expectedUser.Login, expectedUser.Hash))
	ts.NoError(err)
	orderUserID := OrderUserID{OrderNumber: "112233", UserID: "Jhon"}
	_, err = ts.storage.WithRetry(ctx, ts.storage.AddOrder(ctx, orderUserID.OrderNumber, orderUserID.UserID))
	ts.NoError(err)

	testStatuses := []OrderStatusNew{
		{Number: "112233", Status: "PROCESSING", Accrual: 729.98, UploadedAt: time.Now()},
	}

	_, err = ts.storage.WithRetry(ctx, ts.storage.PutStatuses(ctx, &testStatuses))
	ts.NoError(err)

	ordersInterface, err := ts.storage.WithRetry(ctx, ts.storage.GetOrders(ctx, expectedUser.Login))
	ts.NoError(err)
	orders, ok := ordersInterface.([]OrderStatusNew)
	ts.True(ok)
	ts.Equal(testStatuses[0].Number, orders[0].Number)
	ts.Equal(testStatuses[0].Accrual, orders[0].Accrual)
}

func (ts *tSuite) TestGetOrders() {

	ts.T().Log("Тест TestGetOrders()")
	ctx := context.Background()
	ts.TruncateAllTables(ctx)

	expectedUser := User{Login: "Jhon", Hash: "123"}

	_, err := ts.storage.WithRetry(ctx, ts.storage.AddUser(ctx, expectedUser.Login, expectedUser.Hash))
	ts.NoError(err)

	// нет данных
	ordersInterface, err := ts.storage.WithRetry(ctx, ts.storage.GetOrders(ctx, expectedUser.Login))
	ts.NoError(err)
	orders, ok := ordersInterface.([]OrderStatusNew)
	ts.True(ok)
	var orderStatusList []OrderStatusNew
	ts.Equal(orderStatusList, orders, "возврат пустой структуры")

	// добавляем ордер и проверяем
	orderUserID := OrderUserID{OrderNumber: "112233", UserID: "Jhon"}
	_, err = ts.storage.WithRetry(ctx, ts.storage.AddOrder(ctx, orderUserID.OrderNumber, orderUserID.UserID))
	ts.NoError(err)
	ordersInterface, err = ts.storage.WithRetry(ctx, ts.storage.GetOrders(ctx, expectedUser.Login))
	ts.NoError(err)
	orders, ok = ordersInterface.([]OrderStatusNew)
	ts.True(ok)
	ts.Equal("112233", orders[0].Number)
	ts.Equal("NEW", orders[0].Status)

	time.Sleep(time.Second)
	// добавляем еще ордер ордер и проверяем
	orderUserID = OrderUserID{OrderNumber: "1177", UserID: "Jhon"}
	_, err = ts.storage.WithRetry(ctx, ts.storage.AddOrder(ctx, orderUserID.OrderNumber, orderUserID.UserID))
	ts.NoError(err)
	ordersInterface, err = ts.storage.WithRetry(ctx, ts.storage.GetOrders(ctx, expectedUser.Login))
	ts.NoError(err)
	orders, ok = ordersInterface.([]OrderStatusNew)
	ts.True(ok)
	ts.Equal("112233", orders[0].Number)
	ts.Equal("NEW", orders[0].Status)
	ts.Equal("1177", orders[1].Number)
	ts.Equal("NEW", orders[1].Status)
}

func (ts *tSuite) SetupSuite() {
	//база существует
	DatabaseURI := "postgresql://postgres:12345@localhost/gmtest?sslmode=disable"
	ctx := context.Background()
	storage, err := New(ctx, DatabaseURI, "migrations")
	ts.NoError(err)

	ts.storage = storage

}

func (ts *tSuite) TearDownSuite() {

	ts.T().Log("TearDownSuite")
	// if err := m.Down(); err != nil {
	// 	if err.Error() == "no change" {
	// 		println("No migration to rollback")
	// 	} else {
	// 		panic(err)
	// 	}
	// }
	ts.storage.DB.Close()

}

func (ts *tSuite) SetupTest() {

	ts.T().Log("Setup test parameters")

}

func (ts *tSuite) Truncate(ctx context.Context, tableName string) error {

	query := fmt.Sprintf("DELETE FROM %s", tableName)
	_, err := ts.storage.DB.ExecContext(ctx, query)

	if err != nil {
		return fmt.Errorf("ошибка очистки таблицы %s %w", tableName, err)
	}

	return err
}

func (ts *tSuite) TruncateAllTables(ctx context.Context) {

	ts.NoError(ts.Truncate(ctx, "billing"))
	ts.NoError(ts.Truncate(ctx, "orders"))
	ts.NoError(ts.Truncate(ctx, "users"))

}
