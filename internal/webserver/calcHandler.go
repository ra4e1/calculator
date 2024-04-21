package webserver

import (
	"fmt"
	"net/http"

	"github.com/ra4e1/calculator/internal/service"
)

func (web *Webserver) CalcHandler(w http.ResponseWriter, r *http.Request) { //запуск счета и выдача ID
	userId := r.Context().Value("userid").(int64)

	expression := &service.Expression{
		UserId:     userId,
		Ready:      false,
		Err:        nil,
		Value:      nil,
		Expression: r.FormValue("q"),
	}

	requestID, err := web.stateService.AddExpression(expression)
	if err != nil {
		web.ErrorResponse(w, "Ошибка. Не получилось сохранить в базе :(", http.StatusInternalServerError)
		return
	}

	go func(requestID int64) {
		result, err := web.calcService.Calculate(expression.Expression)
		if err != nil {
			expression.Ready = true
			expression.Err = &service.ExpressionError{Msg: err.Error()}
			web.stateService.UpdateExpression(requestID, expression)
			return
		}
		expression.Ready = true
		expression.Value = result
		web.stateService.UpdateExpression(requestID, expression)
	}(requestID)

	fmt.Fprint(w, requestID)
}
