package api

import (
	"log"
	"net/http"

	storage "go-rest-api/storage"
	tasks "go-rest-api/tasks"

	"github.com/gorilla/mux"
)

type APIServer struct {
	store storage.Store
	addr  string
}

func NewAPIServer(addr string, store *storage.Storage) *APIServer {
	return &APIServer{
		store: store,
		addr:  addr,
	}
}

func (server *APIServer) Serve() {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	tasksService := tasks.NewTasksService(server.store)
	tasksHandler := tasks.NewTasksHandler(tasksService)
	tasksHandler.RegisterRoutes(subrouter)

	log.Println("starting the API server at", server.addr)
	log.Fatal(http.ListenAndServe(server.addr, router))
}
