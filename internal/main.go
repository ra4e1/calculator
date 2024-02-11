package main

import "github.com/ra4e1/calculator/internal/application"

func main() { // запуск всего
	config := application.Config{Port: 8081, CalcDelay: 0}
	app := application.NewApplication(config)
	app.Run()
}