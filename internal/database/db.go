package db

import (
	"context"
	"database/sql"
	"fmt"
	"gophermart/pkg/logger"
	"gophermart/utils"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"go.uber.org/zap"
)

var _ StoragerDB = &Storage{}

type StoragerDB interface {
	Close() error
	GetUser(context.Context, string) func(context.Context, *sql.Tx) (interface{}, error)
	AddUser(context.Context, string, string) func(context.Context, *sql.Tx) (interface{}, error)
	GetOrder(context.Context, string) func(context.Context, *sql.Tx) (interface{}, error)
	AddOrder(context.Context, string, string) func(context.Context, *sql.Tx) (interface{}, error)
	GetOrders(context.Context, string) func(context.Context, *sql.Tx) (interface{}, error)
	GetBalance(context.Context, string) func(context.Context, *sql.Tx) (interface{}, error)
	WithdrawBalance(context.Context, string, OrderSum) func(context.Context, *sql.Tx) (interface{}, error)
	WithRetry(context.Context, func(context.Context, *sql.Tx) (interface{}, error)) (interface{}, error)
	GetWithdrawals(context.Context, string) func(context.Context, *sql.Tx) (interface{}, error)
}

type Storage struct {
	DatabaseURI string
	DB          *sql.DB
	logger      *zap.SugaredLogger
}

// подключение к postgress и migrationsUp
func New(ctx context.Context, DatabaseURI string, dbName string) (*Storage, error) {
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

	db, err := migrationsUp(ctx, conn, DatabaseURI, dbName)
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

// func OpenDBConnection(ctx context.Context, DatabaseURI string) (*sql.DB, error) {

// 	db, err := sql.Open("pgx", DatabaseURI)
// 	if err != nil {
// 		return nil, fmt.Errorf("ошибка открытия базы данных %w", err)
// 	}

// 	return db, nil
// }

func (storage *Storage) Close() error {

	return storage.DB.Close()

}

// пока миграции создаем тут. думаю, нужно переписать
func migrationsUp(ctx context.Context, db *sql.DB, DatabaseURI string, dbName string) (*sql.DB, error) {

	// проверяем существует ли база
	var exist string
	row := db.QueryRowContext(ctx, "SELECT datname FROM pg_database where datname=$1;", dbName)
	row.Scan(&exist)

	// создаем если не существует
	if exist != dbName {
		_, err := db.Exec("CREATE DATABASE " + dbName)
		if err != nil {
			return nil, fmt.Errorf("ошибка создания БД %w", err)
		}
	}
	// подключаемся к базе
	db, err := sql.Open("pgx", DatabaseURI+dbName)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия базы данных %w", err)
	}

	// создание таблиц
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (		
		user_id VARCHAR PRIMARY KEY,
		hash VARCHAR NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS orders (
		number VARCHAR PRIMARY KEY,
		user_id VARCHAR NOT NULL,
		uploaded_at timestamp NOT NULL DEFAULT now(),
		FOREIGN KEY (user_id) REFERENCES users(user_id)
	);
	
	CREATE TABLE IF NOT EXISTS billing (
		order_number VARCHAR NOT NULL,
		status VARCHAR NOT NULL,
		accrual int, 
		uploaded_at time NOT NULL,
		time timestamp NOT NULL,
		FOREIGN KEY (order_number) REFERENCES orders(number)
	);	
	
`)
	// #ВОПРОС МЕНТОРУ. Время генерим в таблице orders или передаем извне?

	if err != nil {
		return nil, fmt.Errorf("ошибка при создании таблиц  %w", err)
	}

	return db, nil

}

func (storage *Storage) WithRetry(ctx context.Context, txFunc func(context.Context, *sql.Tx) (interface{}, error)) (interface{}, error) {

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
			if err == nil {
				return nil, fmt.Errorf("ошибка при выполнении commit %w", err)
			}
			break
		}

	}

	return result, nil

}
