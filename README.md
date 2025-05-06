# Biathlon
[![Linter and tests](https://github.com/Arzeeq/biathlon/actions/workflows/main.yml/badge.svg)](https://github.com/Arzeeq/biathlon/actions/workflows/main.yml)

## Description
CLI Application for managing biathlon competitions. See [task](Task.md) for more details

## Run

Run from the root of the project
```bash
go run ./cmd/app/main.go
```

Get manual
```bash
go run ./cmd/app/main.go --help
```

Options
- `--config` - path to config file
- `--events`- path to events file
- `--log` - path to log file
- `--report` - path to report file

## Test
```bash
go test -v ./...
```

