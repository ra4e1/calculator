package application

import (
	"github.com/ra4e1/calculator/internal/service"
	"github.com/ra4e1/calculator/internal/webserver"
)

type Config struct {
	Port      int
	CalcDelay int
	DbName    string
}

type Application struct {
	Cfg Config
	Db  *service.DbService
	Web *webserver.Webserver
}

func NewApplication(config Config) *Application { //создание
	db := service.NewDbService(config.DbName)
	web := webserver.NewWebserver(db)

	return &Application{
		Cfg: config,
		Db:  db,
		Web: web,
	}
}

func (app *Application) Run() int { //запуск
	calc := service.NewCalculatorService(app.Cfg.CalcDelay)
	err := app.Db.Open()
	if err != nil {
		panic(err)
	}
	defer app.Db.Close()

	app.Web.Start(app.Cfg.Port, calc)
	return 0
}
