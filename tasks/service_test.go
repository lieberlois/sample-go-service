package tasks

import (
	"go-rest-api/models"
	"go-rest-api/storage"
	"testing"
	"time"
)

func TestCreateTask(t *testing.T) {
	store := storage.NewMockStore()
	svc := NewTasksService(store)

	t.Run("should create element", func(t *testing.T) {
		// Arrange
		name := "some-task-name"
		task := &models.Task{
			Name: name,
		}

		// Act
		result, err := svc.createTask(task)

		// Assert
		if err != nil {
			t.Error(err)
		}

		if result.Id != 0 {
			t.Errorf("Expected result ID to be %d but got %d", 0, result.Id)
		}

		if result.Name != name {
			t.Errorf("Expected result name to be %s but got %s", name, result.Name)
		}
	})
}

func TestListTasks(t *testing.T) {
	store := storage.NewMockStore()
	svc := NewTasksService(store)

	t.Run("should list one task", func(t *testing.T) {
		// Arrange
		tasks := []*models.Task{
			{
				Id:        1,
				Name:      "hello-world",
				CreatedAt: time.Now(),
			},
		}

		store.Tasks = tasks
		// Act
		result, err := svc.listTasks()

		// Assert
		if err != nil {
			t.Error(err)
		}

		for i, task := range tasks {
			if !task.Equals(result[i]) {
				t.Errorf("Expected task %+v to equal %+v", task, result[i])
			}
		}
	})
}
