package application

import (
	"github.com/ra4e1/calculator/internal/service"
	"github.com/ra4e1/calculator/internal/webserver"
)

type Config struct {
	Port      int
	CalcDelay int
}

type Application struct {
	Cfg Config
	Web *webserver.Webserver
}

func NewApplication(config Config) *Application { //создание
	return &Application{
		Cfg: config,
		Web: webserver.NewWebserver(),
	}
}

func (a *Application) Run() int { //запуск
	calc := service.NewCalculatorService(a.Cfg.CalcDelay)
	a.Web.Start(a.Cfg.Port, calc)
	return 0
}
