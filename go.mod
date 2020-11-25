module github.com/kzmake/example-go-micro-rabbitmq

go 1.14

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/bxcodec/faker/v3 v3.5.0
	github.com/golang/protobuf v1.4.3
	github.com/kzmake/micro-kit v0.4.0
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/broker/rabbitmq/v2 v2.9.1
	google.golang.org/protobuf v1.25.0
)
