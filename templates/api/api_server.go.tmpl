package main

import (
	"fmt"
	"net/http"
)

func (app *application) serverHandlerFunc() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		origin := r.Header.Get("origin")
		method := r.Method
		fmt.Printf("method: %v\n", method)
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Accept", "application/json, multipart/mixed")
		w.Header().Set("Content-Type", "application/json")
		if method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		claims, err := app.auth.ClaimsFromHeader(w, r)
		if err != nil {
			fmt.Printf("auth err: %v\n", err)
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

		fmt.Printf("claims: %v\n", claims)
		app.server.ServeHTTP(w, r)
	})

}
