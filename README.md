# settings-loader
 
## Overview

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
- If the hashes passed through the request and the calculated hash are not equal, the parsed data won't be saved to the database, and the content will be `nil` in the response. The client will receive the new hash and `nil` content for handling the situation.
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

The most useful commands you can find in [Taskfile.yml](Taskfile.yml)
Or you can also run everything using [docker-compose.yml](deploy%2Fdocker-compose.yml).
All e2e tests for non it's only one file [main_test.go](cmd%2Fserver%2Fmain_test.go).

The most useful commands there:

> [!TIP]
> For development better to run db and migration only. After you can run main.go for server or for client side to check
> your code.

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

All configurations places in `/configs/app.yml`
If you need ENV variables you can use format `${DB_HOST:localhost}` - where:

- `DB_HOST` - ENV variable
- `localhost` - default value

**app configurations**:

- `app.hash` - hash function for content. Possible options: `md5`, `sha256`
- `app.path` - folder for json sources

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request. Ensure your code follows the project's
style and includes appropriate tests.

Improvements:

- configurable batches for gorutines (10 for now)
- transaction support for db (at least for aCid)
- fail handling for partial updated data in repository with retry mechanism
- async api + worker task submitting for long data processing
- metrics + traceability
- supporting more targets (repositories)
- supporting more sources (parsers)
- 50% coverage
- code style auto check + sonar (if applicable)
- tests as part of build

## License

This project is licensed under the [MIT License](LICENSE).



