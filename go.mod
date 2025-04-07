module indexer

go 1.24.1

replace example.com/client => ./client

replace example.com/server => ./server

require (
	example.com/client v0.0.0-00010101000000-000000000000
	example.com/server v0.0.0-00010101000000-000000000000
)

require github.com/lib/pq v1.10.9 // indirect
