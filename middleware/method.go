package middleware

import (
	"encoding/json"
	"net/http"
	"todo-list-app/model"
)

func Get(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			res := model.ErrorResponse{
				Error: "Method is not allowed!",
			}
			jsonRes, _ := json.Marshal(res)
			w.WriteHeader(405)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonRes)
			return
		}
		next.ServeHTTP(w, r)
	}) // TODO: replace this
}

func Post(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			res := model.ErrorResponse{
				Error: "Method is not allowed!",
			}
			jsonRes, _ := json.Marshal(res)
			w.WriteHeader(405)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonRes)
			return
		}
		next.ServeHTTP(w, r)
	}) // TODO: replace this
}

func Delete(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			res := model.ErrorResponse{
				Error: "Method is not allowed!",
			}
			jsonRes, _ := json.Marshal(res)
			w.WriteHeader(405)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonRes)
			return
		}
		next.ServeHTTP(w, r)
	}) // TODO: replace this
}
