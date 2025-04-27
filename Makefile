# Include variables from the .envrc file
include .envrc

## help: print this help message
.PHONY: help
help:
	@echo 'Usage':
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'


## run/api: run the cmd/api application
.PHONY: run/api
run:
	go run . -port=${PORT} -env=${ENV}


## tidy: tidy module dependencies and format all .go files
.PHONY: tidy
tidy:
	@echo 'Tidying module dependencies...'
	go mod tidy
	@echo 'Formatting .go files...'
	go fmt ./...

## audit: run quality control checks
.PHONY: audit
audit:
	@echo 'Checking module dependencies...' go mod tidy -diff
	go mod verify
	@echo 'Vetting code...'
	go vet ./...
	go tool staticcheck ./... @echo 'Running tests...'
	go test -race -vet=off ./...

.PHONY: migration
migration:
	@echo 'Creating migration files for ${name}'
	goose -dir ./migrations create ${name} sql -s

.PHONY: build
build:
	@echo 'Building...'
	go build -ldflags='-s' -o=./bin/api .