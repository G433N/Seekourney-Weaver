# Start server

```bash
$ make server
```


The server can now be accessed at http://localhost:8080 or using the client demo. The server's database can be inspected using `$ docker exec -it go-postgres psql -U go-postgres`

# Run tests

```bash
$ make test
```

# Run application or tests on save

## Using `watchexec`

Install `watchexec` using your package manager.

```bash
$ watchexec -e go -r make server
$ watchexec -e go -r make test
```


# Server api

`/all` - Lists all paths in database, probably won't be used in production but
helpful for tests.

`/search` - Query database, will return all paths containing given keywords.
Keywords are sent using http query under the key 'q'.

`/push/paths` - adds one or more paths to the database,
paths are sent using http query under the key 'p'.

`/push/docs` - adds zero or more documents to the database.
Docs are sent using http from an indexer originally dispatched by main server.
Documents are normalized by Core before storage.

`/quit` - Shuts down the server.

# Run client demo

Ensure that the server application is running

```bash
$ go run tui/main.go <command> [<args>]
```

Available commands:

`search` 
`pushpaths`
`pushdocs`
`all`
`index`
`quit`

See server api for more information about corresponding http requests.

# Package structure

```
core - Server and search code, nothing should depend on this package
indexing - Indexers depend on this package
indexer - Indexers for different file types
tui - Tui client for the server
utils - Utilities, should only depend on external packages
```
