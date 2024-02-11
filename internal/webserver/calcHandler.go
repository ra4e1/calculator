package webserver

import (
	"fmt"
	"net/http"

	"github.com/ra4e1/calculator/internal/service"
)

func (web *Webserver) CalcHandler(w http.ResponseWriter, r *http.Request) { //запуск счета и выдача ID
	web.mu.Lock()
	web.requestID++
	requestID := web.requestID
	answer := &service.Answer{Ready: false, Err: nil, Value: nil, Expresion: r.FormValue("q")}
	web.answers[requestID] = answer
	web.stateService.SaveState(web.answers)
	web.mu.Unlock()

	go func() {
		defer web.stateService.SaveState(web.answers)
		result, err := web.calcService.Calculate(answer.Expresion)
		if err != nil {
			answer.Ready = true
			answer.Err = &service.AnswerError{Msg: err.Error()}
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		answer.Ready = true
		answer.Value = result
	}()

	fmt.Fprint(w, requestID)
}
