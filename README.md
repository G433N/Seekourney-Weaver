# Start server

```bash
$ go run core/main.go server
```


The server can now be accessed at http://localhost:8080 or using the client demo. The server's database can be inspected using `$ docker exec -it go-postgres psql -U go-postgres`

# Run tests

```bash
$ go test ./...
```

# Run application or tests on save

## Using `watchexec`

Install `watchexec` using your package manager.

```bash
$ watchexec -e go -r go run core/main.go
$ watchexec -e go -r go test ./...
```


## Server api

`/all` - Lists all paths in database, probably won't be used in production but
helpful for tests

`/search` - Query database, will return all paths containing given keywords.
Keywords are sent using http query under the key 'q'

`/add` - adds one or several paths to the database, paths are sent using http
query under the key 'p'

`/quit` - Shuts down the server

# Run client demo

Ensure that the server application is running

```bash
$ go run tui/main.go <command> [<args>]
```

See server api for more information about possible commands

# Package structure

```
backend - Server and search code, nothing should depend on this package
indexing - Indexers depend on this package
indexer - Indexers for different file types
tui - Tui client for the server
utils - Utilities, should only depend on external packages
```
