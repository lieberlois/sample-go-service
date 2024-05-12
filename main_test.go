package main

import (
	"bytes"
	"context"
	"encoding/json"
	"go-rest-api/models"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

var env = map[string]string{
	"MYSQL_ROOT_PASSWORD": "rootpassword",
	"MYSQL_DATABASE":      "projectmanager",
	"MYSQL_USER":          "user",
	"MYSQL_PASSWORD":      "password",
	"PORT":                "3306",
	"DB_USER":             "user",
	"DB_PASSWORD":         "password",
	"DB_NAME":             "projectmanager",
	"DB_HOST":             "127.0.0.1",
	"JWT_SECRET":          "randomjwtsecret",
}

func mockGetEnv(key string, _ string) string {
	return env[key]
}

func TestWithMySql(t *testing.T) {
	ctx := context.Background()
	mysqlC, err := mysql.RunContainer(ctx, testcontainers.WithImage("mysql:latest"), testcontainers.WithEnv(env))

	defer func() {
		if err := mysqlC.Terminate(ctx); err != nil {
			log.Fatalf("Could not stop mysql: %s", err)
		}
	}()

	if err != nil {
		log.Fatalf("Could not start mysql: %s", err)
	}

	host, err := mysqlC.Endpoint(ctx, "")
	if err != nil {
		log.Fatal(err)
	}
	env["DB_PORT"] = strings.Split(host, ":")[1]

	go run(ctx, mockGetEnv)
	time.Sleep(3 * time.Second)

	t.Run("should return empty list", func(t *testing.T) {
		resp, err := http.Get("http://127.0.0.1:3000/api/v1/tasks")
		if err != nil {
			t.Fatal(err)
		}

		var tasks []models.Task

		data, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(data, &tasks)
		if err != nil {
			t.Fatal(err)
		}

		if len(tasks) > 0 {
			t.Error("Expected empty")
		}
	})

	t.Run("should create task", func(t *testing.T) {
		taskData := struct {
			Name string `json:"name"`
		}{
			Name: "hello",
		}
		buf := &bytes.Buffer{}
		json.NewEncoder(buf).Encode(taskData)

		resp, err := http.Post("http://127.0.0.1:3000/api/v1/tasks", "application/json", buf)
		if err != nil {
			t.Fatal(err)
		}

		var task models.Task

		data, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(data, &task)
		if err != nil {
			t.Fatal(err)
		}

		if task.Name != "hello" {
			t.Error("Wrong task name")
		}

		if task.Id != 1 {
			t.Error("Wrong task id")
		}
	})

	t.Run("should now contain one task", func(t *testing.T) {
		resp, err := http.Get("http://127.0.0.1:3000/api/v1/tasks")
		if err != nil {
			t.Fatal(err)
		}

		var tasks []models.Task

		data, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(data, &tasks)
		if err != nil {
			t.Fatal(err)
		}

		if len(tasks) != 1 {
			t.Error("Expected one task")
		}
	})
}
