module github.com/nnurry/gopds/probabilistics

go 1.22.3

require github.com/nnurry/gopds/protos v0.0.0

require (
	github.com/axiomhq/hyperloglog v0.0.0-20240507144631-af9851f82b27
	github.com/bits-and-blooms/bloom/v3 v3.7.0
	github.com/caarlos0/env/v11 v11.1.0
	github.com/lib/pq v1.10.9
	github.com/redis/go-redis/v9 v9.5.3
	google.golang.org/grpc v1.64.0
)

require (
	github.com/bits-and-blooms/bitset v1.13.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-metro v0.0.0-20180109044635-280f6062b5bc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

replace github.com/nnurry/gopds/protos => ../protos
