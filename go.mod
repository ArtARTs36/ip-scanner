module github.com/artarts36/ip-scanner

go 1.23.3

replace github.com/artarts36/ip-scanner/pkg/ip-scanner-grpc-api => ./pkg/ip-scanner-grpc-api

require (
	github.com/artarts36/go-metrics v0.1.0
	github.com/caarlos0/env/v11 v11.3.1
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.0.1
	github.com/minio/minio-go/v7 v7.0.90
	github.com/oschwald/maxminddb-golang/v2 v2.0.0-beta.2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.19.0
	google.golang.org/grpc v1.68.0
)

require (
	github.com/artarts36/ip-scanner/pkg/ip-scanner-grpc-api v0.0.0-00010101000000-000000000000
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/minio/crc64nvme v1.0.1 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/rs/xid v1.6.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect
	google.golang.org/protobuf v1.35.2 // indirect
)
