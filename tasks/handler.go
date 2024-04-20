package tasks

import (
	"encoding/json"
	"go-rest-api/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Tasks interface {
	createTask(*models.Task) (*models.Task, error)
	listTasks() ([]*models.Task, error)
}

type TasksHandler struct {
	svc Tasks
}

func NewTasksHandler(s Tasks) *TasksHandler {
	return &TasksHandler{
		svc: s,
	}
}

func (s *TasksHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/tasks", s.handleCreateTask).Methods("POST")
	r.HandleFunc("/tasks", s.handleGetTask).Methods("GET")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Hello")) })

}

func (s *TasksHandler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	log.Println("handleCreateTask")
	var task *models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		w.WriteHeader(500)
		return
	}

	task, err := s.svc.createTask(task)
	if err != nil {
		w.WriteHeader(500)
	}

	json.NewEncoder(w).Encode(task)
}

func (s *TasksHandler) handleGetTask(w http.ResponseWriter, r *http.Request) {
	log.Println("handleGetTask")

	tasks, err := s.svc.listTasks()
	if err != nil {
		w.WriteHeader(500)
	}

	json.NewEncoder(w).Encode(tasks)
}
