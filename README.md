# clickhouse-explorer
Simple ClickHouse query executor

[![Build](https://github.com/alxzoomer/clickhouse-explorer/actions/workflows/go.yml/badge.svg)](https://github.com/alxzoomer/clickhouse-explorer/actions/workflows/go.yml)

## Dev environment

### Requirements

- Docker 
- Make
- Go 1.17+

### Docker

Run docker containers from `local-dev/docker-compose/stack.yml`

```shell
make docker/up
```

Stop docker containers

```shell
make docker/down
```

Show logs

```shell
make docker/logs
```

### Init test table

Init simple table with a couple of records

```shell
make docker/data
```

### Run application

Run application

```shell
make run/app
```

### Run tests

Run unit tests

```shell
make test/unit
```
