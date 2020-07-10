module github.com/kzmake/example-go-micro-rabbitmq

go 1.14

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.3

require (
	github.com/kzmake/micro-kit v0.2.0
	github.com/micro/examples v0.2.0
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/broker/rabbitmq/v2 v2.9.1
	github.com/pborman/uuid v1.2.0
)
