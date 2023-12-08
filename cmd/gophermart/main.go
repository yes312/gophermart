package main

import (
	"context"
	"flag"
	"fmt"

	"gophermart/internal/app"
	"gophermart/internal/config"

	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var f config.Flags

// func init() {

// #ВопросМентору: стоит ли строчки ниже спрятать в функцию в пакет config  или в отдельный пакет?
// flag.StringVar(&f.A, "a", "localhost:8080", "IP adress")
// flag.StringVar(&f.D, "d", "postgres://postgres:12345@localhost:5432/", "database uri")
// flag.StringVar(&f.R, "r", "http://127.0.0.1:8081", "ACCRUAL_SYSTEM_ADDRESS")
// flag.StringVar(&f.A, "a", "", "IP adress")
// flag.StringVar(&f.D, "d", "", "database uri")
// flag.StringVar(&f.R, "r", "", "ACCRUAL_SYSTEM_ADDRESS")

// }

func main() {

	flag.StringVar(&f.A, "a", "localhost:8080", "IP adress")
	flag.StringVar(&f.D, "d", "postgres://postgres:12345@localhost:5432/gm?sslmode=disable", "database uri")
	flag.StringVar(&f.R, "r", "http://127.0.0.1:8081", "ACCRUAL_SYSTEM_ADDRESS")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Восстановлено от паники MAIN:", r)
		}
	}()

	log.Println("====Запуск MAIN====")
	flag.Parse()

	config, err := config.NewConfig(f)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("====Запуск после конфига====")
	// #ВопросМентору: нужно ли graceful shutdown реализовывать как отдельную функцию или метод и нужен ли для этого отдельный пакет?
	// --------------------
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {

		<-c
		cancel()
		fmt.Println("Завершение по сигналу с клавиатуры. ")
		os.Exit(0)

	}()
	// -------------------

	s := app.New(ctx, config)
	wg := &sync.WaitGroup{}
	defer func() {
		wg.Wait()
		if err := s.Close(); err != nil {
			log.Println("ошибка при закрытии сервера:", err)
		} else {
			log.Println("работа сервера успешно завершена")
		}

	}()

	if err := s.Start(ctx, wg); err != nil {
		log.Fatal(err)
	}

}
