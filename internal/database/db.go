package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gophermart/pkg/logger"
	"gophermart/utils"
	"path/filepath"
	"runtime"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var _ StoragerDB = &Storage{}

type dbOperation func(context.Context, *sql.Tx) (interface{}, error)

type StoragerDB interface {
	Close() error
	GetUser(context.Context, string) dbOperation
	AddUser(context.Context, string, string) dbOperation
	GetOrder(context.Context, string) dbOperation
	AddOrder(context.Context, string, string) dbOperation
	GetOrders(context.Context, string) dbOperation
	GetBalance(context.Context, string) dbOperation
	WithdrawBalance(context.Context, string, OrderSum) dbOperation
	WithRetry(context.Context, dbOperation) (interface{}, error)
	GetWithdrawals(context.Context, string) dbOperation
	GetNewProcessedOrders(context.Context) dbOperation
	PutStatuses(context.Context, *[]OrderStatus) dbOperation
}

type Storage struct {
	DatabaseURI string
	DB          *sql.DB
	logger      *zap.SugaredLogger
}

func getProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}

// подключение к postgress и migrationsUp
func New(ctx context.Context, DatabaseURI string, MigrationsPath string) (*Storage, error) {
	// #ВопросМентору объявить новый логгер или передать его с параметром функции из app.go
	logger, err := logger.NewLogger("Info")
	if err != nil {
		return nil, err
	}
	// подключаемся к postgres
	conn, err := sql.Open("pgx", DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия базы данных %w", err)
	}

	fmt.Println("Project root: ", getProjectRoot())

	// filename := "gophermart/cmd/gophermart/main.go"
	// dir := filepath.Dir(filename)
	// fmt.Println(dir)
	// absoluteFilePath, err := filepath.Abs("")
	// if err != nil {
	// 	return nil, err
	// }

	// absoluteFilePath := fmt.Sprint(getProjectRoot(), "\\", MigrationsPath)

	db, err := migrationsUp(ctx, conn, DatabaseURI, MigrationsPath)
	if err != nil {
		return nil, err
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ошибка открытия базы данных(Ping) %w", err)
	}

	return &Storage{
		DatabaseURI: DatabaseURI,
		DB:          db,
		logger:      logger,
	}, nil
}

func (storage *Storage) Close() error {

	return storage.DB.Close()

}

// пока миграции создаем тут. думаю, нужно переписать
func migrationsUp(ctx context.Context, db *sql.DB, DatabaseURI string, migrations string) (*sql.DB, error) {

	path := fmt.Sprintf("file://%s", migrations)
	m, err := migrate.New(path, DatabaseURI)
	if err != nil {
		return nil, err
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}
	return db, nil

	// проверяем существует ли база
	// var exist string
	// row := db.QueryRowContext(ctx, "SELECT datname FROM pg_database where datname=$1;", dbName)
	// row.Scan(&exist)

	// // создаем если не существует
	// if exist != dbName {
	// 	_, err := db.Exec("CREATE DATABASE " + dbName)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("ошибка создания БД %w", err)
	// 	}
	// }
	// // подключаемся к базе
	// db, err := sql.Open("pgx", DatabaseURI+dbName)
	// if err != nil {
	// 	return nil, fmt.Errorf("ошибка открытия базы данных %w", err)
	// }

	// создание таблиц
	// _, err := db.Exec(`
	// 	CREATE TABLE IF NOT EXISTS users (
	// 		user_id VARCHAR PRIMARY KEY,
	// 		hash VARCHAR NOT NULL
	// 	);

	// 	CREATE TABLE IF NOT EXISTS orders (
	// 		number VARCHAR PRIMARY KEY,
	// 		user_id VARCHAR NOT NULL,
	// 		uploaded_at timestamp NOT NULL,
	// 		FOREIGN KEY (user_id) REFERENCES users(user_id)
	// 	);

	// 	CREATE TABLE IF NOT EXISTS billing (
	// 		order_number VARCHAR NOT NULL,
	// 		status VARCHAR NOT NULL,
	// 		accrual int,
	// 		uploaded_at timestamp NOT NULL,
	// 		time timestamp NOT NULL,
	// 		FOREIGN KEY (order_number) REFERENCES orders(number),
	// 		CONSTRAINT unique_order_number_status UNIQUE (order_number, status)
	// 	);

	// `)

	// if err != nil {
	// 	return nil, fmt.Errorf("ошибка при создании таблиц  %w", err)
	// }

	// return db, nil

}

func (storage *Storage) WithRetry(ctx context.Context, txFunc dbOperation) (interface{}, error) {

	var result interface{}
	pauseDurations := []int{0, 1, 3, 5}

	for _, pause := range pauseDurations {

		select {
		case <-ctx.Done():
			return nil, nil
		case <-time.After(time.Duration(pause) * time.Second):
		}

		tx, err := storage.DB.Begin()
		defer tx.Rollback()
		if err != nil {
			return nil, fmt.Errorf("ошибка при создании транзакции %w", err)
		}

		result, err = txFunc(ctx, tx)

		if err != nil {
			if !utils.OnDialErr(err) {
				return nil, fmt.Errorf("НЕвостановимая ошибка %w", err)
			}
			storage.logger.Info("восстановимая ошибка %v", err)
		} else {
			err = tx.Commit()
			if err != nil {
				return nil, fmt.Errorf("ошибка при выполнении commit %w", err)
			}
			break
		}

	}

	return result, nil

}
