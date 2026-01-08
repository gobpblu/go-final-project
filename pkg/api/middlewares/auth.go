package middlewares

import (
	"net/http"

	"go-final-project/pkg/utils"
)

func AuthMiddleware(password string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(password) > 0 {
				var jwt string

				cookie, err := r.Cookie("token")
				if err == nil {
					jwt = cookie.Value
				}

				valid := utils.IsTokenValid(jwt, password)

				if !valid {
					utils.WriteAuthenticationRequiredError(w)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
