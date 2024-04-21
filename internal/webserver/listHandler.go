package webserver

import (
	"fmt"
	"net/http"
)

func (web *Webserver) ListHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userid").(int64)
	//fmt.Fprintf(w, "USER ID: %d\n", userId)

	expressions, err := web.stateService.FindUserAllExpressions(userId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	for _, exp := range expressions {
		if !exp.Ready {
			fmt.Fprintf(w, "%d) %s - решается\n", exp.ID, exp.Expression)
			continue
		}
		if exp.Err != nil {
			fmt.Fprintf(w, "%d) %s - ошибка, %s\n", exp.ID, exp.Expression, exp.Err.Msg)
			continue
		}
		if exp.Ready {
			fmt.Fprintf(w, "%d) %s - решено\n", exp.ID, exp.Expression)
			continue
		}
	}
}
