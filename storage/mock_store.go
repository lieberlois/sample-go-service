package storage

import (
	"go-rest-api/models"
)

type MockStore struct {
	tasks []*models.Task
}

func NewMockStore() *MockStore {
	return &MockStore{
		tasks: make([]*models.Task, 0),
	}
}

func (s *MockStore) CreateTask(task *models.Task) (*models.Task, error) {
	id := int64(len(s.tasks))
	task.Id = id

	s.tasks = append(s.tasks, task)
	return task, nil
}

func (s *MockStore) ListTasks() ([]*models.Task, error) {
	return s.tasks, nil
}
