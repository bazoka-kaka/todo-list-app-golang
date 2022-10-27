package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"
	"todo-list-app/db"
	"todo-list-app/middleware"
	"todo-list-app/model"

	"github.com/google/uuid"
)

func Register(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./view/register.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./view/login.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	username := fmt.Sprintf("%s", r.Context().Value("username"))
	tasks := db.Task[username]
	data := map[string]interface{}{
		"username": username,
		"tasks":    tasks,
	}
	tmpl, err := template.ParseFiles("./view/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	reenterPassword := r.FormValue("reenter-password")
	if password != reenterPassword {
		middleware.ShowMessage(w, "Password Doesn't Match!", 400)
		return
	}
	db.Users[username] = password
	db.Task[username] = []model.Todo{}
	middleware.ShowMessage(w, "Register Success!", 200)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	expectedPwd, ok := db.Users[username]
	if !ok || expectedPwd != password {
		middleware.ShowMessage(w, "Wrong Password or Username!", 401)
		return
	}
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(30 * time.Minute)
	db.Sessions[sessionToken] = model.Session{
		Username: username,
		Expiry:   expiresAt,
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Path:    "/",
		Expires: expiresAt,
	})

	middleware.ShowMessage(w, "Login Successful!", 200)
}

func main() {
	// without auth
	http.Handle("/register", middleware.Get(http.HandlerFunc(Register)))
	http.Handle("/login", middleware.Get(http.HandlerFunc(Login)))

	http.Handle("/user/register", middleware.Post(http.HandlerFunc(HandleRegister)))
	http.Handle("/user/login", middleware.Post(http.HandlerFunc(HandleLogin)))

	// using auth
	http.Handle("/", middleware.Auth(middleware.Get(http.HandlerFunc(Home))))

	fmt.Println("server running on port 3000")
	http.ListenAndServe(":3000", nil)
}
