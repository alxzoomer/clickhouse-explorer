
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
