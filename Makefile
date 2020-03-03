.PHONY: build
build:
	go build -v ./cmd/apiserver

setup:
	go get -u github.com/BurntSushi/toml
	go get -u github.com/sirupsen/logrus
	go get -u github.com/gorilla/mux
	go get -u github.com/stretchr/testify
	go get -u github.com/go-sql-driver/mysql
	go get -u golang.org/x/crypto
	go get -u github.com/go-ozzo/ozzo-validation/v4
	go get -u github.com/go-ozzo/ozzo-validation/v4/is
	go get -u github.com/gorilla/sessions
	go get -u github.com/gorilla/handlers
	go get -u github.com/google/uuid

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build
