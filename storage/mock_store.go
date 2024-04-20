package storage

import (
	"go-rest-api/models"
)

type MockStore struct {
	Tasks []*models.Task
}

func NewMockStore() *MockStore {
	return &MockStore{
		Tasks: make([]*models.Task, 0),
	}
}

func (s *MockStore) CreateTask(task *models.Task) (*models.Task, error) {
	id := int64(len(s.Tasks))
	task.Id = id

	s.Tasks = append(s.Tasks, task)
	return task, nil
}

func (s *MockStore) ListTasks() ([]*models.Task, error) {
	return s.Tasks, nil
}
