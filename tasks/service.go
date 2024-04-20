package tasks

import (
	"go-rest-api/models"

	storage "go-rest-api/storage"
)

type TasksService struct {
	store storage.Store
}

func NewTasksService(s storage.Store) *TasksService {
	return &TasksService{
		store: s,
	}
}

func (s *TasksService) createTask(task *models.Task) (*models.Task, error) {
	result, err := s.store.CreateTask(task)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *TasksService) listTasks() ([]*models.Task, error) {
	result, err := s.store.ListTasks()
	if err != nil {
		return nil, err
	}

	return result, nil
}
