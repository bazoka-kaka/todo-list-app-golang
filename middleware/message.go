package middleware

import (
	"html/template"
	"net/http"
)

func ShowMessage(w http.ResponseWriter, msg string, code int) {
	message := map[string]string{
		"message": msg,
	}
	tmpl, err := template.ParseFiles("./view/message.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	w.WriteHeader(code)
	if err := tmpl.Execute(w, message); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
	}
}
