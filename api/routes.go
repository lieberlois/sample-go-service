package api

import (
	tasks "go-rest-api/tasks"

	"github.com/gorilla/mux"
)

func addRoutes(
	mux *mux.Router,
	tasksService *tasks.TasksService,
) {
	mux.Handle("/tasks", tasks.HandleCreateTask(tasksService)).Methods("POST")
	mux.Handle("/tasks", tasks.HandleGetTask(tasksService)).Methods("GET")
}
