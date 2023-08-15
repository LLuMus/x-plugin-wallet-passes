package auth

import "net/http"

const cookieName = "wallet-passes-auth"

func CookieIntoHeader() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookies := r.Cookies()
			authHeader := r.Header.Get("Authorization")

			if len(cookies) > 0 && authHeader == "" {
				for _, cookie := range cookies {
					if cookie.Name == cookieName {
						r.Header.Set("Authorization", "Bearer "+cookie.Value)
						break
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
