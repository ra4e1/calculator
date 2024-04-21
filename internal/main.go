package main

import (
	"os"
	"strconv"

	"github.com/ra4e1/calculator/internal/application"
)

func main() { // запуск всего
	port, err := strconv.Atoi(os.Getenv("CALC_APP_PORT"))
	if err != nil {
		port = 8081
	}
	calcDelay, err := strconv.Atoi(os.Getenv("CALC_APP_DELAY"))
	if err != nil {
		calcDelay = 10
	}

	config := application.Config{
		Port:      port,       // Порт http-сервера
		CalcDelay: calcDelay,  // Задержка подсчета в секундах
		DbName:    "store.db", // имя файла базы данных
	}
	app := application.NewApplication(config)
	app.Run()
}
