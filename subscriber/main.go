package main

import (
	"context"

	"github.com/kzmake/micro-kit/pkg/logger/technical"
	plogger "github.com/kzmake/micro-kit/pkg/wrapper/logger"
	proto "github.com/micro/examples/pubsub/srv/proto"
	micro "github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/util/log"
	rabbitmq "github.com/micro/go-plugins/broker/rabbitmq/v2"
)

const topic = "kzmake.queue.example"

func subscribe(ctx context.Context, event *proto.Event) error {
	md, _ := metadata.FromContext(ctx)
	technical.Infof("[Subscribe %s] Received event %+v with metadata %+v\n", topic, event, md)

	return nil
}

func main() {
	b := rabbitmq.NewBroker()

	service := micro.NewService(
		micro.Broker(b),
		micro.WrapSubscriber(plogger.NewSubscriberWrapper()),
	)
	service.Init()

	subopts := broker.NewSubscribeOptions(
		broker.Queue("queue.default"),
		broker.DisableAutoAck(),
		rabbitmq.DurableQueue(),
	)

	// register
	micro.RegisterSubscriber(topic, service.Server(), subscribe, server.SubscriberContext(subopts.Context), server.SubscriberQueue("queue.example"))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
