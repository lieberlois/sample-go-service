package storage

import (
	"go-rest-api/models"
)

type Store interface {
	CreateTask(*models.Task) (*models.Task, error)
	ListTasks() ([]*models.Task, error)
}
