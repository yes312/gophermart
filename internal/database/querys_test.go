package db

// import (
// 	"context"
// 	"fmt"
// 	"testing"

// 	"github.com/stretchr/testify/suite"
// )

// type tSuite struct {
// 	suite.Suite
// 	storage *Storage
// }

// const testDBName = "gmtest"

// func TestSuiteTest(t *testing.T) {
// 	suite.Run(t, &tSuite{})
// }

// func (ts *tSuite) TestAddGetUser() {

// 	ts.T().Log("Тест AddUser() и GetUser()")
// 	ctx := context.Background()
// 	ts.NoError(ts.Truncate(ctx))

// 	user := User{Login: "Jhon", Hash: "123"}
// 	err := ts.storage.AddUser(ctx, user.Login, user.Hash)
// 	ts.NoError(err)

// 	expectedUser, err := ts.storage.GetUser(ctx, user.Login)
// 	ts.NoError(err)
// 	ts.Equal(expectedUser.Login, user.Login)
// 	ts.Equal(expectedUser.Hash, user.Hash)

// }
// func (ts *tSuite) SetupSuite() {

// 	DatabaseURI := "postgres://postgres:12345@localhost:5432/"

// 	ctx := context.Background()

// 	storage, err := New(ctx, DatabaseURI, testDBName)
// 	ts.NoError(err)

// 	ts.storage = storage

// }

// func (ts *tSuite) TearDownSuite() {

// 	ts.T().Log("TearDownSuite")
// 	ts.storage.DB.Close()

// }

// func (ts *tSuite) SetupTest() {

// 	ts.T().Log("Setup test parameters")

// }

// func (ts *tSuite) Truncate(ctx context.Context) error {

// 	_, err := ts.storage.DB.ExecContext(ctx, `TRUNCATE TABLE users;`)
// 	if err != nil {
// 		return fmt.Errorf("ошибка создания таблицы metrics %w", err)
// 	}

// 	return nil
// }
