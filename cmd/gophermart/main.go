package main

import (
	"context"
	"flag"

	"gophermart/internal/app"
	"gophermart/internal/config"
	"gophermart/pkg/logger"

	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var f config.Flags

func init() {

	// #ВопросМентору: стоит ли строчки ниже спрятать в функцию в пакет config  или в отдельный пакет?

	flag.StringVar(&f.A, "a", "localhost:8081", "IP adress")
	// flag.StringVar(&f.D, "d", "postgresql://postgres:12345@localhost/gmtest?sslmode=disable", "database uri")
	flag.StringVar(&f.R, "r", "http://127.0.0.1:8080", "ACCRUAL_SYSTEM_ADDRESS")
	flag.StringVar(&f.D, "d", "postgresql://postgres:12345@localhost/gmtest?sslmode=disable", "database uri")
}

func main() {

	flag.Parse()

	config, err := config.NewConfig(f)
	if err != nil {
		log.Fatal(err)
	}
	logger, err := logger.NewLogger(config.LoggerLevel)

	if err != nil {
		logger.Error(err)
		os.Exit(0)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {

		<-c
		cancel()
		logger.Info("Завершение по сигналу с клавиатуры. ")
		os.Exit(0)

	}()

	s := app.New(ctx, config)
	wg := &sync.WaitGroup{}
	defer func() {
		wg.Wait()
		if err := s.Close(); err != nil {
			logger.Info("ошибка при закрытии сервера:", err)
		} else {
			logger.Info("работа сервера успешно завершена")
		}

	}()

	if err := s.Start(ctx, logger, wg); err != nil {
		logger.Error(err)
		os.Exit(0)
	}

}
