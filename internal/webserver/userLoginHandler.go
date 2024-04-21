package webserver

import (
	"encoding/json"
	"errors"
	"net/http"
)

type UserLoginData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (web *Webserver) UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginData UserLoginData
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&loginData)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			web.ErrorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			web.ErrorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}

	token, err := web.userService.Login(loginData.Login, loginData.Password)
	if err != nil {
		web.ErrorResponse(w, "Login error: "+err.Error(), http.StatusForbidden)
		return
	}

	resp := make(map[string]string)
	resp["result"] = "OK"
	resp["token"] = token
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
