package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/kzmake/micro-kit/pkg/logger/technical"
	plogger "github.com/kzmake/micro-kit/pkg/wrapper/logger"
	micro "github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/metadata"
	rabbitmq "github.com/micro/go-plugins/broker/rabbitmq/v2"

	"github.com/kzmake/example-go-micro-rabbitmq/proto"
)

const topic = "test.hoge"

var actions = []string{
	"CreateHoge",
	"DeleteHoge",
	"UpdateHoge",
}

func newPublishOptions(opts ...broker.PublishOption) broker.PublishOptions {
	opt := broker.PublishOptions{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func newEvent(id string) *proto.HogeEvent {
	return &proto.HogeEvent{
		Id:          id,
		TimestampNs: time.Now().UnixNano(),
		Action:      actions[rand.Intn(len(actions))], // rand action
		Data: map[string]string{
			"username": faker.Username(),
			"email":    faker.Email(),
		},
	}
}

func publish(p micro.Publisher, e *proto.HogeEvent, opts ...client.PublishOption) {
	md := metadata.Metadata{
		"x-request-id": faker.UUIDHyphenated(),
		"client-ip":    faker.IPv4(),
	}

	// publish an event
	ctx := metadata.NewContext(context.Background(), md)
	if err := p.Publish(ctx, e, opts...); err != nil {
		technical.Infof("event の Publish に失敗しました: %+v", err)
	} else {
		technical.Infof("%s: %+v with metadata %+v", e.GetId(), e, md)
	}
}

func main() {
	service := micro.NewService(
		micro.Broker(rabbitmq.NewBroker(
			rabbitmq.ExchangeName("test.hoge"),
			rabbitmq.PrefetchCount(1),
		)),
		micro.WrapSubscriber(plogger.NewSubscriberWrapper()),
	)
	service.Init()

	// create publisher
	p := micro.NewPublisher(topic, service.Client())

	publish(p, newEvent("1st: (priority 1)"), client.PublishContext(newPublishOptions(rabbitmq.Priority(1)).Context))
	publish(p, newEvent("2nd: (priority 7)"), client.PublishContext(newPublishOptions(rabbitmq.Priority(7)).Context))
	publish(p, newEvent("3rd: (priority 3)"), client.PublishContext(newPublishOptions(rabbitmq.Priority(3)).Context))
	publish(p, newEvent("4th: (priority 5)"), client.PublishContext(newPublishOptions(rabbitmq.Priority(5)).Context))
	publish(p, newEvent("5th: (priority 9)"), client.PublishContext(newPublishOptions(rabbitmq.Priority(9)).Context))
}
