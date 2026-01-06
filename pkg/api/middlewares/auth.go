package middlewares

import (
	"fmt"
	"go-final-project/pkg/utils"
	"net/http"
	"os"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		fmt.Println("ENV: ", pass)
		if len(pass) > 0 {
			var jwt string // JWT-токен из куки

			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}

			valid := utils.IsTokenValid(jwt)

			fmt.Println("VALID: ", valid)

			if !valid {
				utils.WriteAuthenticationRequiredError(w)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
