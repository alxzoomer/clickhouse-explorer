
.PHONY: help
## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## run/app: run application
run/app:
	go run ./cmd/clickhouse-explorer run

## test/unit: run unit tests
test/unit:
	go test ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

## vet: examines Go source code and reports suspicious constructs
.PHONY: vet
vet:
	go vet ./...

## fmt: format Go code
.PHONY: fmt
fmt:
	go fmt ./...

## pre-commit: run pre-commit actions to format and check the code
.PHONY: precommit
pre-commit: fmt vet

compose_project=clickhouse-explorer
## docker/up: run docker containers
.PHONY: docker/up
docker/up:
	docker-compose -p ${compose_project} -f ./local-dev/docker-compose/stack.yml up -d

## docker/down: stop docker containers
.PHONY: docker/down
docker/down:
	docker-compose -p ${compose_project} -explorer -f ./local-dev/docker-compose/stack.yml down

## docker/logs: read docker containers log output
.PHONY: docker/logs
docker/logs:
	docker-compose -p ${compose_project} -f ./local-dev/docker-compose/stack.yml logs -f

## db/init: create test database, example_table and fill few records
.PHONY: db/init
db/init:
	cat local-dev/example-table.sql | \
		docker-compose -p ${compose_project} -f ./local-dev/docker-compose/stack.yml exec -T \
		-- clickhouse clickhouse-client -mn --verbose
