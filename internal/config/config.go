package config

import (
	"gophermart/utils"
	"log"
	"os"
	"time"
)

// Сервис должен поддерживать конфигурирование следующими методами:
// адрес и порт запуска сервиса: переменная окружения ОС RUN_ADDRESS или флаг -a;
// адрес подключения к базе данных: переменная окружения ОС DATABASE_URI или флаг -d;
// адрес системы расчёта начислений: переменная окружения ОС ACCRUAL_SYSTEM_ADDRESS или флаг -r.

type Flags struct {
	A string // RUN_ADDRESS
	D string // DATABASE_URI
	R string // ACCRUAL_SYSTEM_ADDRESS
}

type Config struct {
	RunAdress                string
	AccrualSysremAdress      string
	AccrualRequestInterval   int
	AccuralPuttingDBInterval int
	DatabaseURI              string
	LoggerLevel              string
	Key                      string
	TokenExp                 time.Duration
}

func NewConfig(flag Flags) (*Config, error) {
	log.Println("NewConfig=================")

	c := Config{}
	if buf, ok := os.LookupEnv("RUN_ADDRESS"); ok {
		c.RunAdress = buf
	} else {
		var err error
		if c.RunAdress, err = utils.GetValidURL(flag.A); err != nil {
			return &Config{}, utils.ErrorWrongURLFlag
		}
	}

	if buf, ok := os.LookupEnv("DATABASE_URI"); ok {
		c.DatabaseURI = buf
	} else {
		c.DatabaseURI = flag.D
	}

	if buf, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); ok {
		c.AccrualSysremAdress = buf
	} else {
		c.AccrualSysremAdress = flag.R
	}

	c.LoggerLevel = "Info"

	// secretKeyPass := "configs/key.toml"
	// _, err := toml.DecodeFile(secretKeyPass, &c)
	// if err != nil {
	// 	return &Config{}, err
	// }
	c.Key = "12345"
	c.TokenExp = time.Hour * 999
	c.AccrualRequestInterval = 10
	c.AccuralPuttingDBInterval = 5

	return &c, nil

}
