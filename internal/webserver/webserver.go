package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/ra4e1/calculator/internal/service"
)

type Webserver struct {
	calcService  *service.CalculatorService
	stateService *service.DbStateService
	userService  *service.UserService
}

func NewWebserver(db *service.DbService) *Webserver { //создание
	return &Webserver{
		stateService: service.NewDbStateService(db),
		userService:  service.NewUserService(db),
	}
}

func (w *Webserver) loadState() {
}

func (web *Webserver) ErrorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func (w *Webserver) createRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	calc := http.HandlerFunc(w.CalcHandler)
	mux.Handle("/calc", w.authMiddleware(calc))

	answer := http.HandlerFunc(w.AnswerHandler)
	mux.Handle("/answer", w.authMiddleware(answer))

	list := http.HandlerFunc(w.ListHandler)
	mux.Handle("/list", w.authMiddleware(list))

	userRegister := http.HandlerFunc(w.UserRegisterHandler)
	mux.Handle("/api/v1/register", userRegister)

	userLogin := http.HandlerFunc(w.UserLoginHandler)
	mux.Handle("/api/v1/login", userLogin)

	return mux
}

func (w *Webserver) Start(port int, calcService *service.CalculatorService) { //запуск
	w.calcService = calcService
	w.loadState()
	mux := w.createRoutes()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
		if err != nil {
			fmt.Println("Ошибка при запуске http сервера", err)
		}
	}()
	fmt.Printf("Запуск веб-сервера http://localhost:%d\n", port)
	wg.Wait()
}
