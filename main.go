package main

import (
	"calculator/application"
)

func main() { // запуск всего
	config := application.Config{Port: 8081, CalcDelay: 0}
	app := application.NewApplication(config)
	app.Run()
}
