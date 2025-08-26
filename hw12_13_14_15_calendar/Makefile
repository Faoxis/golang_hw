CALENDAR_BIN := "./bin/calendar"
SCHEDULER_BIN := "./bin/calendar_scheduler"
SENDER_BIN := "./bin/calendar_sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(CALENDAR_BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(SCHEDULER_BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
	go build -v -o $(SENDER_BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

run: build
	$(CALENDAR_BIN) -config ./configs/config.yml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(CALENDAR_BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.57.2

LINT_BIN := $(shell go env GOPATH)/bin/golangci-lint
lint: install-lint-deps
	$(LINT_BIN) run ./...

lint-docker:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.57.2 golangci-lint run

generate:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/EventService.proto

# Docker commands
up:
	cd env && docker-compose up -d

down:
	cd env && docker-compose down

logs:
	cd env && docker-compose logs -f

# Integration tests
integration-tests:
	cd env && docker-compose -f docker-compose.yml -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from integration-tests

# Kubernetes commands
k8s-apply:
	kubectl apply -f k8s/

k8s-delete:
	kubectl delete -f k8s/

k8s-status:
	kubectl get all -n calendar

# Helm commands
helm-install:
	helm install calendar ./helm/calendar

helm-upgrade:
	helm upgrade calendar ./helm/calendar

helm-uninstall:
	helm uninstall calendar

helm-status:
	helm status calendar

.PHONY: build run build-img run-img version test lint lint-docker generate up down logs integration-tests k8s-apply k8s-delete k8s-status helm-install helm-upgrade helm-uninstall helm-status
