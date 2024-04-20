package tasks

import (
	"bytes"
	"encoding/json"
	"go-rest-api/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

type MockTasksServiceProxy struct {
	create func(*models.Task) (*models.Task, error)
	list   func() ([]*models.Task, error)
}

func (mts *MockTasksServiceProxy) createTask(t *models.Task) (*models.Task, error) {
	return mts.create(t)
}

func (mts *MockTasksServiceProxy) listTasks() ([]*models.Task, error) {
	return mts.list()
}

func TestHandleListTask(t *testing.T) {
	svc := new(MockTasksServiceProxy)
	handler := NewTasksHandler(svc)

	sampleModelList := []*models.Task{
		{
			Id:        0,
			Name:      "Hello World",
			CreatedAt: time.Now(),
		},
	}

	data := []struct {
		name     string
		list     func() ([]*models.Task, error)
		expected []*models.Task
	}{
		{
			name: "empty list",
			list: func() ([]*models.Task, error) {
				return sampleModelList, nil
			},
			expected: sampleModelList,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			// Arrange
			svc.list = d.list

			req, err := http.NewRequest("GET", "/tasks", &bytes.Buffer{})
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/tasks", handler.handleGetTask).Methods("GET")

			// Act
			router.ServeHTTP(rr, req)

			// Assert
			if rr.Code != http.StatusOK {
				t.Errorf("Expected HTTP status %d but got %d instead", http.StatusOK, rr.Code)
			}

			// Check the response body is what we expect
			actual := make([]*models.Task, 0)
			err = json.NewDecoder(rr.Body).Decode(&actual)

			if err != nil {
				t.Error(err)
			}

			if len(actual) != len(d.expected) {
				t.Fatalf("Expected %d tasks, got %d tasks", len(d.expected), len(actual))
			}

			for i, task := range actual {
				if !task.Equals(d.expected[i]) {
					t.Errorf("Expected task at index %d to be %+v, got %+v", i, d.expected[i], task)
				}
			}
		})
	}
}

func TestHandleCreateTask(t *testing.T) {
	svc := new(MockTasksServiceProxy)
	handler := NewTasksHandler(svc)

	t.Run("should return task with id", func(t *testing.T) {
		// Arrange
		name := "sample-name"
		expectedId := 5

		svc.create = func(t *models.Task) (*models.Task, error) {
			t.Id = int64(expectedId)
			return t, nil
		}

		task := &models.Task{
			Name: name,
		}

		body, _ := json.Marshal(task)
		req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/tasks", handler.handleCreateTask).Methods("POST")

		// Act
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusOK {
			t.Errorf("Expected HTTP status %d but got %d instead", http.StatusOK, rr.Code)
		}

		// Check the response body is what we expect
		actual := new(models.Task)
		err = json.NewDecoder(rr.Body).Decode(actual)

		if err != nil {
			t.Error(err)
		}

		if actual.Id != int64(expectedId) {
			t.Errorf("expected ID to be %d but got %d", expectedId, actual.Id)
		}

		if actual.Name != name {
			t.Errorf("expected name to be %s but got %s", name, actual.Name)
		}
	})
}
