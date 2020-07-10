package main

import (
	"context"
	"fmt"
	"time"

	"github.com/kzmake/micro-kit/pkg/logger/technical"
	plogger "github.com/kzmake/micro-kit/pkg/wrapper/logger"
	proto "github.com/micro/examples/pubsub/srv/proto"
	micro "github.com/micro/go-micro/v2"
	rabbitmq "github.com/micro/go-plugins/broker/rabbitmq/v2"
	"github.com/pborman/uuid"
)

const topic = "kzmake.queue.example"

func publish(topic string, p micro.Publisher) {
	t := time.NewTicker(time.Second)

	for range t.C {
		e := &proto.Event{
			Id:        uuid.NewUUID().String(),
			Timestamp: time.Now().Unix(),
			Message:   fmt.Sprintf("Messaging you all day on %s", topic),
		}

		// publish an event
		if err := p.Publish(context.Background(), e); err != nil {
			technical.Infof("error publishing: %v", err)
		} else {
			technical.Infof("[Publish %s] Received event %v with metadata", topic, e)
		}
	}
}

func main() {
	service := micro.NewService(
		micro.Broker(rabbitmq.NewBroker()),
		micro.WrapSubscriber(plogger.NewSubscriberWrapper()),
	)
	service.Init()

	// create publisher
	pub := micro.NewPublisher(topic, service.Client())

	// publish
	go publish(topic, pub)

	// block forever
	select {}
}
