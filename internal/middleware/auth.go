package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler, jwtSecret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[AUTH] Checking authentication for %s %s", r.Method, r.URL.Path)

		cookie, err := r.Cookie("jwt_token")
		if err != nil {
			log.Printf("[AUTH] No JWT cookie found: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Printf("[AUTH] JWT cookie found, validating token")
		tokenString := cookie.Value
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("[AUTH] Invalid signing method: %v", token.Method)
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			log.Printf("[AUTH] Token parsing failed: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			log.Printf("[AUTH] Token is invalid")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Printf("[AUTH] Token validated successfully")
		next.ServeHTTP(w, r)
	})
}

func getCookieToken(r *http.Request) string {
	cookie, err := r.Cookie("jwt_token")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func getAuthToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return getCookieToken(r)
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return getCookieToken(r)
	}

	return parts[1]
}
