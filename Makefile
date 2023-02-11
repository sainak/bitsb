.DEFAULT_GOAL := help
.PHONY: build bin make migrations

# --- Tooling & Variables ----------------------------------------------------------------
# Exporting bin folder to the path for makefile
export PATH   := $(PWD)/bin:$(PATH)
# Default Shell
export SHELL  := bash
# Type of OS: Linux or Darwin.
export OSTYPE := $(shell uname -s)

# export .env
ifneq (,$(wildcard ./.env))
	include .env
	export
endif

include ./make/tools.Makefile
include ./make/help.Makefile

COMPOSE := docker-compose -f ./compose.yaml

VERSION ?= $(shell git describe --tags --always --dirty)


# ~~~ Development Environment ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

install-deps: migrate air gotestsum tparse mockery ## Install Development Dependencies (localy).
deps: $(MIGRATE) $(AIR) $(GOTESTSUM) $(TPARSE) $(MOCKERY) ## Checks for Global Development Dependencies.
deps:
	@echo "Required Tools Are Available"

dev: dev-env dev-air ## Startup / Spinup Docker Compose and air for development.

up: ## Run application (within a Docker-Compose)
	@if [ ! -f .env ]; then \
		echo "Creating new .env file"; \
		cp .env.example .env; \
	fi
	@ $(COMPOSE) up --wait -d

logs: ## Show logs (within a Docker-Compose)
	@ $(COMPOSE) logs -f --no-log-prefix api; true

dev-env: ## Bootstrap Environment (with a Docker-Compose).
	@ $(COMPOSE) up -d --build db

dev-env-test: dev-env ## Run application (within a Docker-Compose)
	#@ $(MAKE) image-build
	@ $(COMPOSE) up --build api

dev-air: $(AIR) ## Starts AIR ( Continuous Development app).
	@ air

db-shell: ## Run database shell (In local environment)
	@ psql "${DB_DSN}"

docker-db-shell: ## Run database shell (within a Docker-Compose)
	@ $(COMPOSE) exec db psql -U postgres

down: ## Stop docker services
	@ $(COMPOSE) down

teardown: ## Teardown (removes volumes, tmp files, etc...)
	@ $(COMPOSE) down --remove-orphans -v


# ~~~ Docker Build ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

image-build: ## Build Docker Image
	@ echo "Docker Build"
	@ DOCKER_BUILDKIT=0 docker build \
		--file deployments/Dockerfile \
		--build-arg VERSION=${VERSION} \
		--tag bitsb \
			.


# ~~~ Code Actions ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

TESTS_ARGS := --format testname --jsonfile gotestsum.json.out
TESTS_ARGS += --max-fails 2
TESTS_ARGS += -- ./...
TESTS_ARGS += -test.failfast
TESTS_ARGS += -test.parallel $(shell nproc)
TESTS_ARGS += -test.count 1
TESTS_ARGS += -test.timeout 5s
TESTS_ARGS += -coverpkg ./...
TESTS_ARGS += -covermode=atomic
TESTS_ARGS += -coverprofile coverage.out
TESTS_ARGS += -race

run-tests: $(GOTESTSUM)
	@ gotestsum $(TESTS_ARGS) -short
	@ cat coverage.out | grep -v "mocks" > coverage.out.tmp
	@ mv coverage.out.tmp coverage.out

tests: run-tests $(TPARSE) ## Run Tests & parse details
	@cat gotestsum.json.out | tparse -all -notests

lint: $(GOLANGCI) ## Runs golangci-lint with predefined configuration
	@echo "Applying linter"
	golangci-lint version
	golangci-lint run -c .golangci.yaml ./...

build: ## Builds binary
	@ printf "Building aplication... "
	@ go build \
		-trimpath  \
		-ldflags "-X 'main.version=${VERSION}'"\
		-o ./build/engine \
		./app/
	@ echo "done"

build-race: ## Builds binary (with -race flag)
	@ printf "Building aplication with race flag... "
	@ go build \
		-trimpath  \
		-race      \
		-ldflags "-X 'main.version=${VERSION}-race'"\
		-o ./build/engine \
		./app/
	@ echo "done"

mocks: $(MOCKERY) ## Runs go generte ./...
	@ printf "Generating mocks..."
	@ mockery --all --output ./domain/mocks


# ~~~ Database Migrations ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

migrate-up: $(MIGRATE) ## Apply all (or N up) migrations.
	@ read -p "How many migration you wants to perform (default value: [all]): " N; \
	migrate  -database "${DB_DSN}" -path=migrations up ${NN}

migrate-down: $(MIGRATE) ## Apply all (or N down) migrations.
	@ read -p "How many migration you wants to perform (default value: [all]): " N; \
	migrate  -database "${DB_DSN}" -path=migrations down ${NN}

migrate-drop: $(MIGRATE) ## Drop everything inside the database.
	migrate  -database "${DB_DSN}" -path=migrations drop

migrations: $(MIGRATE) ## Create a set of up/down migrations with a specified name.
	@ read -p "Please provide name for the migration: " Name; \
	migrate create -ext sql -dir migrations $${Name}


# ~~~ Cleanup ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

clean: clean-artifacts clean-docker ## Cleanup (removes volumes, tmp files, etc...)

clean-artifacts: ## Removes Artifacts (*.out)
	@printf "Cleanning artifacts... "
	@rm -f *.out
	@rm -r build
	@echo "done."

clean-docker: ## Tear down docker containers and images
	@ $(MAKE) docker-teardown
	@ docker image prune -f
