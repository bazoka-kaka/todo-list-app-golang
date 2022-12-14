package main

import (
	"fmt"
	"net/http"
	"os"
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

	tasks := []map[string]string{}
	for _, item := range db.Task[username] {
		task := map[string]string{
			"task": item.Task,
		}
		if item.Done {
			task["done"] = "true"
		} else {
			task["done"] = "false"
		}
		tasks = append(tasks, task)
	}
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

func AddTask(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./view/addTask.html")
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

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	username := fmt.Sprintf("%s", r.Context().Value("username"))

	tasks := []map[string]string{}
	for _, item := range db.Task[username] {
		task := map[string]string{
			"task": item.Task,
		}
		if item.Done {
			task["done"] = "true"
		} else {
			task["done"] = "false"
		}
		tasks = append(tasks, task)
	}
	data := map[string]interface{}{
		"username": username,
		"tasks":    tasks,
	}

	tmpl, err := template.ParseFiles("./view/deleteTask.html")
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

	for key := range db.Users {
		if key == username {
			middleware.ShowMessage(w, "User Already Exist!", 400)
			return
		}
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

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		middleware.ShowMessage(w, "You are not Logged in!", 401)
		return
	}
	sessionToken := c.Value

	delete(db.Sessions, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Path:    "/",
		Expires: time.Now(),
	})

	middleware.ShowMessage(w, "Logout Successful!", 200)
}

func HandleAddTask(w http.ResponseWriter, r *http.Request) {
	username := fmt.Sprintf("%s", r.Context().Value("username"))

	for _, item := range db.Task[username] {
		if item.Task == r.FormValue("task") {
			middleware.ShowMessage(w, "Task Already Exists!", 400)
			return
		}
	}

	var todo model.Todo

	todo.Id = uuid.NewString()
	todo.Task = r.FormValue("task")
	if r.FormValue("done") == "on" {
		todo.Done = true
	} else {
		todo.Done = false
	}

	db.Task[username] = append(db.Task[username], todo)
	middleware.ShowMessage(w, "Task Added!", 201)
}

func HandleUpdateTask(w http.ResponseWriter, r *http.Request) {
	username := fmt.Sprintf("%s", r.Context().Value("username"))

	// update tasks
	for i, item := range db.Task[username] {
		if r.FormValue(item.Task) == "on" {
			db.Task[username][i].Done = true
		} else {
			db.Task[username][i].Done = false
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	username := fmt.Sprintf("%s", r.Context().Value("username"))

	newTasks := []model.Todo{}
	deleteExists := false
	for _, item := range db.Task[username] {
		if r.FormValue(item.Task) != "on" {
			newTasks = append(newTasks, item)
		} else {
			deleteExists = true
		}
	}
	db.Task[username] = newTasks
	if deleteExists {
		middleware.ShowMessage(w, "Task(s) Deleted!", 200)
	} else {
		middleware.ShowMessage(w, "No Task(s) Deleted", 200)
	}
}

func main() {
	// without auth
	http.Handle("/register", middleware.Get(http.HandlerFunc(Register)))
	http.Handle("/login", middleware.Get(http.HandlerFunc(Login)))

	http.Handle("/user/register", middleware.Post(http.HandlerFunc(HandleRegister)))
	http.Handle("/user/login", middleware.Post(http.HandlerFunc(HandleLogin)))
	http.Handle("/user/logout", middleware.Get(http.HandlerFunc(HandleLogout)))

	// using auth
	http.Handle("/", middleware.Auth(middleware.Get(http.HandlerFunc(Home))))
	http.Handle("/task/add", middleware.Auth(middleware.Get(http.HandlerFunc(AddTask))))
	http.Handle("/task/delete", middleware.Auth(middleware.Get(http.HandlerFunc(DeleteTask))))

	http.Handle("/task/handler/add", middleware.Auth(middleware.Post(http.HandlerFunc(HandleAddTask))))
	http.Handle("/task/handler/update", middleware.Auth(middleware.Get(http.HandlerFunc(HandleUpdateTask))))
	http.Handle("/task/handler/delete", middleware.Auth(middleware.Post(http.HandlerFunc(HandleDeleteTask))))

	PORT := os.Getenv("PORT")
	// PORT := "3000"
	fmt.Println("server running on port " + PORT)
	http.ListenAndServe(":"+PORT, nil)
}
