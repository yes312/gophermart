package db

import (
	"context"
	"database/sql"
	"fmt"
	"gophermart/internal/services"
	"gophermart/pkg/logger"

	_ "github.com/jackc/pgx/v5/stdlib"

	"go.uber.org/zap"
)

const dbName = "gm"

var _ StoragerDB = &Storage{}

type StoragerDB interface {
	Close() error
	GetUser(services.UserAuthInfo)
}

type Storage struct {
	DatabaseURI string
	DB          *sql.DB
	logger      *zap.SugaredLogger
}

// подключение к postgress и migrationsUp
func New(ctx context.Context, DatabaseURI string) (*Storage, error) {
	// #ВопросМентору объявить новый логгер или передать его с параметром функции из app.go
	logger, err := logger.NewLogger("Info")
	if err != nil {
		return nil, err
	}
	// подключаемся к postgres
	db, err := sql.Open("pgx", DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия базы данных %w", err)
	}

	err = migrationsUp(ctx, db, DatabaseURI)
	if err != nil {
		return nil, err
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ошибка открытия базы данных(Ping) %w", err)
	}
	defer db.Close()

	return &Storage{
		DatabaseURI: DatabaseURI,
		DB:          db,
		logger:      logger,
	}, nil
}

func OpenDBConnection(ctx context.Context, DatabaseURI string) (*sql.DB, error) {

	db, err := sql.Open("pgx", DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия базы данных %w", err)
	}

	return db, nil
}

func (storage *Storage) Close() error {

	return storage.DB.Close()

}

// пока миграции создаем тут. думаю, нужно переписать
func migrationsUp(ctx context.Context, db *sql.DB, DatabaseURI string) error {

	// проверяем существует ли база
	var exist string
	row := db.QueryRowContext(ctx, "SELECT datname FROM pg_database where datname=$1;", dbName)
	row.Scan(&exist)

	// создаем если не существует
	if exist != dbName {
		_, err := db.Exec("CREATE DATABASE " + dbName)
		if err != nil {
			return fmt.Errorf("ошибка создания БД %w", err)
		}
	}
	// подключаемся к базе
	db, err := sql.Open("pgx", DatabaseURI+dbName)
	if err != nil {
		return fmt.Errorf("ошибка открытия базы данных %w", err)
	}

	// нужно для использования типа uuid
	_, err = db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if err != nil {
		return fmt.Errorf("ошибка uuid-ossp %w", err)
	}

	// создание таблицы
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		uid uuid DEFAULT uuid_generate_v4 (),
		name VARCHAR NOT NULL,
		hash VARCHAR NOT NULL
	)
`)
	if err != nil {
		return fmt.Errorf("ошибка при создании таблицы users %w", err)
	}

	return nil

}
