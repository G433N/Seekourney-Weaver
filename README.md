
# Run application
```bash
$ go run main.go

```

# Run tests
```bash
$ go test ./...
```

# Run application or tests on save

## Using `watchexec`
Install `watchexec` using your package manager.
```bash
$ watchexec -e go -r go run main.go
$ watchexec -e go -r go test ./...
```
