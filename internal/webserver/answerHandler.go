package webserver

import (
	"fmt"
	"net/http"
	"strconv"
)

func (web *Webserver) AnswerHandler(w http.ResponseWriter, r *http.Request) { //выдача ответа
	userId := r.Context().Value("userid").(int64)
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	expression, err := web.stateService.FindUserExpression(userId, int64(id))
	if err != nil {
		web.ErrorResponse(w, err.Error(), http.StatusNotFound)
		return
	}

	if !expression.Ready {
		fmt.Fprintf(w, "ожидайте решения")
		return
	}
	if expression.Err != nil {
		fmt.Fprintln(w, expression.Err.Msg)
		return
	}

	fmt.Fprint(w, expression.Value)
}
