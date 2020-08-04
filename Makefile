SHELL := /bin/bash
BASEDIR = $(shell pwd)

export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=off

APP_NAME=appboot
APP_VERSION=1.0.0
IMAGE_PREFIX=appboot/${APP_NAME}
IMAGE_NAME=${IMAGE_PREFIX}:v${APP_VERSION}
IMAGE_JENKINSAPI=${IMAGE_PREFIX}:jenkinsapi-v${APP_VERSION}
IMAGE_LATEST=${IMAGE_PREFIX}:latest
IMAGE_JENKINSAPI_LATEST=${IMAGE_PREFIX}:jenkinsapi-latest

fmt:
	gofmt -w .
mod:
	go mod tidy
lint:
	golangci-lint run
.PHONY: build
build:
	go build -o appboot cmd/appboot/main.go
	go build -o server cmd/server/main.go
build-docker:
	sh build/package/build.sh ${IMAGE_NAME}
.PHONY: test
test:
	go test -gcflags=-l -coverpkg=./... -coverprofile=coverage.data ./...
.PHONY: web
web:
	cd web/appboot; \
	npm run serve
help:
	@echo "fmt - gofmt"
	@echo "mod - go mod tidy"
	@echo "lint - run golangci-lint"
	@echo "utest - unit test"