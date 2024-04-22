# settings-loader

## Overview

Settings Loader is a Go project designed to manage and load various application settings from different sources (you can
mount your volume). For now, it's limited only JSON files sources and only one db table support, but it could be
extended.

This project is template for the future microservices. Current layout based
on https://github.com/golang-standards/project-layout

### API

Api of this microservice - RPC Functions

#### LoaderController.LoadComponent

Main rpc function. Current implementation:

- The function reads file from the volume/disc (`path=<rootPath>/<type>/<version>.json`).
- Saves parsed data to the postgresql database, table `settings`.
- Calculates file content hash (by default uses `MD5Hash` function)
- If hashes, passed through request and calculated are not equal, then parsed data won't save to db and content will be
  nil in response. Client will get a new hash and nil content for handling the situation.
- If file doesn't exist, then error returns.
- Use defaults if they are not present in the payload.

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

This project is licensed under the MIT License.



