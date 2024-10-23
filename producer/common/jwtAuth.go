package common

import (
	"context"
	"net/http"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
)

const ROLE_WRITER = "writer"

// secret encryption key
var jwtKey = []byte(GetEnv().JwtKey)

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// generate a jwt token for a user and role
func GenerateJwtTokenString(username string, role string) (string, error) {
	// uncomment if you need the key to expire (e.g. 30 days)
	// expirationTime := jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour))
	claims := &Claims{
		Username: username,
		Role:     role,
		// RegisteredClaims: jwt.RegisteredClaims{
		// 	ExpiresAt: expirationTime,
		// },
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Middleware to authenticate
func MwAuthenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		// parse and validate token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// store user in the request context for later use
		r = r.WithContext(context.WithValue(r.Context(), "username", claims.Username))
		r = r.WithContext(context.WithValue(r.Context(), "role", claims.Role))

		next.ServeHTTP(w, r)
	})
}

// Middleware to authorise "writer" role
func MwAuthorizeWriter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value("role")
		if role != ROLE_WRITER {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
