package api

import (
	storage "go-rest-api/storage"
	tasks "go-rest-api/tasks"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAPIServer(store *storage.Storage) http.Handler {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	tasksService := tasks.NewTasksService(store)
	addRoutes(subrouter, tasksService)

	return router
}
