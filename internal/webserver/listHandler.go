package webserver

import (
	"fmt"
	"net/http"
	"sort"
)

func (web *Webserver) ListHandler(w http.ResponseWriter, r *http.Request) {
	keys := make([]int, 0, len(web.answers))
	for k := range web.answers {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, requestID := range keys {
		answer := web.answers[requestID]
		if !answer.Ready {
			fmt.Fprintf(w, "%d) %s - решается\n", requestID, answer.Expresion)
			continue
		}
		if answer.Err != nil {
			fmt.Fprintf(w, "%d) %s - ошибка, %s\n", requestID, answer.Expresion, answer.Err.Msg)
			continue
		}
		if answer.Ready {
			fmt.Fprintf(w, "%d) %s - решено\n", requestID, answer.Expresion)
			continue
		}
	}
}
