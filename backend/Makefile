.PHONY: tidy
tidy:
	go mod tidy -v
	# go fmt ./...

.PHONY: build
build:
	go build -o=/tmp/bin/traintrack ./cmd/api
	
.PHONY: run
run: build
	/tmp/bin/traintrack

## migrations/new name=$1: create a new database migration
.PHONY: migrations/new
migrations/new:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest create -seq -ext=.sql -dir=./assets/migrations ${name}

## migrations/up: apply all up database migrations
.PHONY: migrations/up
migrations/up:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="${DB_DSN}" up

## migrations/down: apply all down database migrations
.PHONY: migrations/down
migrations/down:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="${DB_DSN}" down

## migrations/goto version=$1: migrate to a specific version number
.PHONY: migrations/goto
migrations/goto:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="${DB_DSN}" goto ${version}

## migrations/force version=$1: force database migration
.PHONY: migrations/force
migrations/force:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="${DB_DSN}" force ${version}

## migrations/version: print the current in-use migration version
.PHONY: migrations/version
migrations/version:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="${DB_DSN}" version
