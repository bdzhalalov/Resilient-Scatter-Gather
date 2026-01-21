build:
	docker-compose build

run:
	docker-compose up -d
	sleep 1
	docker-compose ps

start: build run
	@echo "All services are up"

TEST_CONTAINER_NAME := resilient-scatter-gather-test

.PHONY: test build run

# Запуск тестов
test:
	docker build --target builder -t $(TEST_CONTAINER_NAME) .
	docker run --rm $(TEST_CONTAINER_NAME) go test -v ./internal/handler

stop:
	docker-compose stop