run: build
	@./bin/api

.PHONY: build
build:
	@go build -o ./bin/api

.PHONY: test
test:
	@go test -v ./...

.PHONY: get 
get:
	curl localhost:3000/api/v1/tasks

.PHONY: post 
post:
	curl -X POST localhost:3000/api/v1/tasks -d '{"name": "Hello world"}'