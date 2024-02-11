package webserver

import (
	"fmt"
	"net/http"
	"strconv"
)

func (web *Webserver) AnswerHandler(w http.ResponseWriter, r *http.Request) { //выдача ответа
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	answer, ok := web.answers[id]
	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if !answer.Ready {
		fmt.Fprintf(w, "ожидайте решения")
		return
	}
	if answer.Err != nil {
		fmt.Fprintln(w, answer.Err.Msg)
		return
	}

	fmt.Fprint(w, answer.Value)
}
