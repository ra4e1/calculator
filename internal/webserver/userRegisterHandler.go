package webserver

import (
	"encoding/json"
	"errors"
	"net/http"
)

type UserRegData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (web *Webserver) UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var regData UserRegData
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&regData)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			web.ErrorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			web.ErrorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}

	_, err = web.userService.Register(regData.Login, regData.Password)
	if err != nil {
		web.ErrorResponse(w, "App error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	web.ErrorResponse(w, "OK", http.StatusOK)
}
