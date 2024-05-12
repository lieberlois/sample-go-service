package tasks

import (
	"encoding/json"
	"go-rest-api/models"
	"go-rest-api/util"
	"log"
	"net/http"
)

type TasksSvc interface {
	createTask(*models.Task) (*models.Task, error)
	listTasks() ([]*models.Task, error)
}

func HandleCreateTask(svc TasksSvc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("handleCreateTask")
		var task *models.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			w.WriteHeader(500)
			return
		}

		task, err := svc.createTask(task)
		if err != nil {
			w.WriteHeader(500)
		}

		err = util.Encode(w, r, http.StatusCreated, task)
		if err != nil {
			w.WriteHeader(500)
		}
	})
}

func HandleGetTask(svc TasksSvc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tasks, err := svc.listTasks()
		if err != nil {
			w.WriteHeader(500)
		}

		err = util.Encode(w, r, http.StatusOK, tasks)
		if err != nil {
			w.WriteHeader(500)
		}
	})
}
