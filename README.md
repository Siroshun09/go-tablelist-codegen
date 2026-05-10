# go-tablelist-codegen

[![Release](https://img.shields.io/github/release/Siroshun09/go-tablelist-codegen)](https://github.com/Siroshun09/go-tablelist-codegen/releases/latest)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/Siroshun09/go-tablelist-codegen/ci.yml?branch=main)
![GitHub](https://img.shields.io/github/license/Siroshun09/go-tablelist-codegen)

A Go library to generate structs that define table and column names in the database.

## Usage

### CLI

- MySQL: `go run github.com/Siroshun09/go-tablelist-codegen/cmd/mysql@<version> <flags>`
- PostgreSQL: `go run github.com/Siroshun09/go-tablelist-codegen/cmd/postgresql@<version> <flags>`
- SQLite: `go run github.com/Siroshun09/go-tablelist-codegen/cmd/sqlite@<version> <flags>`

### Library

```shell
go get github.com/Siroshun09/go-tablelist-codegen
```

## Flags

- `--package-name <package name>`
- `--output <output filepath>` (Optional)
  - If not specified, the codegen prints the generated code to stdout
- `--debug <true/false>` (Optional)
  - Prints stacktrace on error
  - Default: `false`

### For MySQL and PostgreSQL

- `--host <host>`
- `--port <port>`
- `--user <username>`
- `--password <password>`
- `--database <database name>`
- `--sslmode <sslmode>` (Optional)
  - Default: `disable`

### For SQLite

- `--schema-file <schema filepath>`

## License

This project is under the Apache License version 2.0. Please see [LICENSE](LICENSE) for more info.

Copyright © 2025-2026, Siroshun09
