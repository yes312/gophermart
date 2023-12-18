package config

import (
	"fmt"
	"gophermart/utils"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
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
	MigrationsPath           string
	NumberOfWorkers          int
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
	// ВОПРОСМЕНТОРУ в этом месте я пытаюсь прочитать конфиг из файла key.toml в windows файл прочитается, если папка с файлом будет лежать в одной
	// папке с main. В linux на git тестах. прочитается если папка с файлом будет лежать в корне git репозитория.
	// как решить этот вопрос?
	filepathStr := filepath.Join("configs", "key.toml")
	absPath, err := filepath.Abs(filepathStr)

	// _, err = os.Stat(filepathStr)
	// C:\Users\Andrey\go\src\github.com\yes312\gophermart\configs
	if err != nil {
		fmt.Println("Ошибка при получении абсолютного пути:", err)
		return &Config{}, err
	}
	// _, err = os.Stat(absPath)

	// if err == nil {
	// 	fmt.Println("Файл", absPath, "существует.")
	// } else if os.IsNotExist(err) {
	// 	fmt.Println("Файла", absPath, "не существует.")
	// } else {
	// 	fmt.Println("Ошибка при проверке файла:", err)
	// }
	// filePath := filepath.Join(absPath, "configs", "key.toml")
	_, err = toml.DecodeFile(absPath, &c)
	if err != nil {
		// c.Key = "secret"
		return &Config{}, err
	}
	// filePath = filepath.Join("configs", "config.toml")
	// _, err = toml.DecodeFile(filePath, &c)
	// if err != nil {
	// 	return &Config{}, err
	// }
	// c.TokenExp = c.TokenExp * time.Hour

	c.MigrationsPath = "migrations"
	c.TokenExp = time.Hour * 999
	c.AccrualRequestInterval = 1
	c.AccuralPuttingDBInterval = 1
	c.NumberOfWorkers = 3
	return &c, nil

}
