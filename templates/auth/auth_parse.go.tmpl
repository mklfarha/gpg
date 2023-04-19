package auth

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func ParseToken(tknStr string, w *http.ResponseWriter) (*Claims, error) {
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			if w != nil {
				(*w).WriteHeader(http.StatusUnauthorized)
			}
			return nil, err
		}
		if w != nil {
			(*w).WriteHeader(http.StatusBadRequest)
		}
		return nil, err
	}
	if !tkn.Valid {
		if w != nil {
			(*w).WriteHeader(http.StatusUnauthorized)
		}
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
