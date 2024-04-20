package storage

import (
	"database/sql"
	"go-rest-api/models"
)

type Storage struct {
	db *sql.DB
}

func NewDbStorage(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateTask(task *models.Task) (*models.Task, error) {
	rows, err := s.db.Exec("INSERT INTO tasks (name) VALUES (?)", task.Name)
	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}

	task.Id = id
	return task, nil
}

func (s *Storage) ListTasks() ([]*models.Task, error) {
	rows, err := s.db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}

	tasks := make([]*models.Task, 0)

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		task := new(models.Task)
		if err := rows.Scan(&task.Id, &task.Name, &task.CreatedAt); err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return tasks, err
	}
	return tasks, nil
}
