# settings-loader

## Overview

[![Go Report Card](https://goreportcard.com/badge/github.com/ngirchev/settings-loader?style=flat-square)](https://goreportcard.com/report/github.com/ngirchev/settings-loader)
[![Go Reference](https://pkg.go.dev/badge/github.com/ngirchev/settings-loader.svg)](https://pkg.go.dev/github.com/ngirchev/settings-loader)

Settings Loader is a Go project designed to manage and load various application settings from different sources (you can
mount your own volume). Currently, it supports only JSON file sources and one database table, but it could be
extended in the future.

This project serves as a template for future microservices. The current layout is based
on https://github.com/golang-standards/project-layout.

### API

The API of this microservice comprises RPC functions.

#### LoaderController.LoadComponent

The main RPC function, with the following current implementation:

- Reads a file from the volume/disk (`path=<rootPath>/<type>/<version>.json`).
- Saves the parsed data to the PostgreSQL database, table `settings`.
- Calculates a file content hash (by default, it uses the `MD5Hash` function).
- If the hashes passed through the request and the calculated hash are not equal, the parsed data won't be saved to the
  database, and the content will be `nil` in the response. The client will receive the new hash and `nil` content for
  handling the situation.
- If the file doesn't exist, it returns an error.
- Uses default values if they are not present in the payload.

**Input Payload:**

```
type Request struct {
	Type    string // default core
	Version string // default 1.0.0
	Hash    []byte // expected hash
}
```

**Response**

```
type Response struct {
	Type    string // same as in request
	Version string // same as in request
	Hash    []byte // new hash
	Content []byte // content is nil if 'expected hash' != 'new hash'
}
```

### Structure

- `cmd/server/main.go` - run server side
- `cmd/client/main.go` - run client rpc code for development
- `cmd/generator/main.go` - run for print random json data for test
- `configs/app.yml` - app configuration
- `deploy/*` - docker related files
- `internal/*` - all business logic code
- `resources/*` - json files, sources for parsing
- `schema/*` - migration files

## Build and Run

You can find the most useful commands in [Taskfile.yml](Taskfile.yml). Alternatively, you can run everything
using [docker-compose.yml](deploy/docker-compose.yml). The end-to-end tests are located in one
file: [main_test.go](cmd/server/main_test.go).

Here are some useful commands to get started:

> [!TIP]
> For development, it's best to run only the database and migration. Afterward, you can run `main.go` to start the
> server or the client, depending on what you need to test.

**DOWN ALL**

``` shell
docker-compose -f deploy/docker-compose.yml down -v
```

**UP ALL**

``` shell
docker-compose -f deploy/docker-compose.yml up -d
```

**UP DB**

``` shell
docker-compose -f deploy/docker-compose.yml up -d migrate
```

## Configuration

All configurations are placed in `/configs/app.yml`.
If you need to use environment variables, you can follow this format: `${DB_HOST:localhost}`. Here's what the parts
represent:

- `DB_HOST`: the environment variable
- `localhost`: the default value if the environment variable isn't set

**Application Configurations**:

- `app.hash`: the hash function for content. Possible options are `md5` or `sha256`.
- `app.path`: the folder for JSON sources.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request. Ensure that your code follows the
project's style and includes appropriate tests.

Improvements or TODO:

- Check concurrent access to the database and write tests for it.
- Check concurrent access to the file system and write tests.
- Configurable batch sizes for goroutines (currently set to 10).
- Transaction support for the database (at least for ACID compliance).
- Failure handling for partially updated data in the repository with a retry mechanism.
- Asynchronous API with worker task submission for long data processing.
- Metrics and traceability.
- Support for more targets (repositories).
- Support for more sources (parsers).
- 50% code coverage.
- Automated code style checks and Sonar integration (if applicable).
- Tests as part of the build process.

## License

This project is licensed under the [MIT License](LICENSE).



