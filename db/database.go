package db

import (
	"todo-list-app/model"
)

var Users = map[string]string{}
var Task = map[string][]model.Todo{}

var Sessions = map[string]model.Session{}
