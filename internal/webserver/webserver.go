package webserver

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/ra4e1/calculator/internal/service"
)

type Webserver struct {
	requestID    int
	answers      map[int]*service.Answer
	mu           sync.Mutex
	calcService  *service.CalculatorService
	stateService *service.StateService
}

func NewWebserver() *Webserver { //создание
	return &Webserver{
		requestID:    0,
		answers:      make(map[int]*service.Answer),
		mu:           sync.Mutex{},
		stateService: service.NewStateService(),
	}
}

func (w *Webserver) loadState() {
	answers, requestID, err := w.stateService.RestoreState()
	if err == nil {
		w.answers = answers
		w.requestID = requestID
	} else {
		fmt.Println(err)
	}
}

func (w *Webserver) createRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	calc := http.HandlerFunc(w.CalcHandler)
	mux.Handle("/calc", calc)

	answer := http.HandlerFunc(w.AnswerHandler)
	mux.Handle("/answer", answer)

	list := http.HandlerFunc(w.ListHandler)
	mux.Handle("/list", list)
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
