package webserver

import (
	"context"
	"net/http"

	"github.com/ra4e1/calculator/internal/service"
)

type RequestUserContext struct {
	UserId int64
}

func (web *Webserver) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := ""
		if r.Header["Token"] != nil {
			token = r.Header["Token"][0]
		} else {
			token = r.FormValue("token")
		}

		if token == "" {
			web.ErrorResponse(w, "Auth error", http.StatusUnauthorized)
			return
		}

		authServie := service.NewAuthService()
		userId, err := authServie.ValidateToken(token)
		if err != nil {
			web.ErrorResponse(w, "Auth error: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userid", userId)

		// Передаем управление следующему обработчику
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
