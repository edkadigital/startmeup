ENV ?= dev

DATABASE_URL ?= postgres://startmeup:secret@localhost:5432/startmeup?sslmode=disable
TEST_DATABASE_URL ?= postgres://startmeup:secret@localhost:5432/startmeup_test?sslmode=disable

.PHONY: help
help: ## Print make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: ent-install
ent-install: ## Install Ent code-generation module
	go get entgo.io/ent/cmd/ent

.PHONY: air-install
air-install: ## Install air
	go install github.com/air-verse/air@latest

.PHONY: ent-gen
ent-gen: ## Generate Ent code
	go generate ./ent

.PHONY: ent-new
ent-new: ## Create a new Ent entity (ie, make ent-new name=MyEntity)
	go run entgo.io/ent/cmd/ent new $(name)

.PHONY: admin
admin: ## Create a new admin user (ie, make admin email=myemail@web.com)
	go run cmd/admin/main.go --email=$(email)

.PHONY: migrate
migrate: ## Run database migrations including River queue tables
ifeq ($(ENV),test)
	@echo "Running migrations on TEST database ($(TEST_DATABASE_URL))..."
	DATABASE_URL=$(TEST_DATABASE_URL) go run cmd/migrate/main.go
else
	@echo "Running migrations on default database ($(DATABASE_URL))..."
	DATABASE_URL=$(DATABASE_URL) go run cmd/migrate/main.go
endif

.PHONY: build
build: ## Build all applications (web, worker, migrate)
	go build -o bin/web cmd/web/main.go
	go build -o bin/worker cmd/worker/main.go
	go build -o bin/migrate cmd/migrate/main.go

.PHONY: run
run: ## Run the web application only
	clear
	go run cmd/web/main.go

.PHONY: run-worker
run-worker: ## Run the background task worker
	clear
	go run cmd/worker/main.go

.PHONY: watch
watch: ## Run the application and watch for changes with air to automatically rebuild
	clear
	air

.PHONY: test
test: ## Run all tests
ifeq ($(ENV),test)
	@echo "Running tests on TEST database ($(TEST_DATABASE_URL))..."
	DATABASE_URL=$(TEST_DATABASE_URL) go test ./...
else
	@echo "Running tests on default database ($(DATABASE_URL))..."
	DATABASE_URL=$(DATABASE_URL) go test ./...
endif

.PHONY: check-updates
check-updates: ## Check for direct dependency updates
	go list -u -m -f '{{if not .Indirect}}{{.}}{{end}}' all | grep "\["

.PHONY: fmt
fmt: ## Format the code
	go fmt ./...

.PHONY: lint
lint: ## Lint the code
	go vet ./...