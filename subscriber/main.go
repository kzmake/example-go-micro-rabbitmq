package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/kzmake/micro-kit/pkg/logger"
	"github.com/kzmake/micro-kit/pkg/logger/technical"
	plogger "github.com/kzmake/micro-kit/pkg/wrapper/logger"
	micro "github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/server"
	rabbitmq "github.com/micro/go-plugins/broker/rabbitmq/v2"

	"github.com/kzmake/example-go-micro-rabbitmq/proto"
)

const (
	exchange = "test.hoge.exchange"
	topic    = "test.hoge.create"
	queue    = "test.hoge.queue"
)

func subscribe(ctx context.Context, e *proto.HogeEvent) error {
	// recv
	log.Printf("%s\n", e.GetId())
	// log.Printf("\t%+v with metadata %+v\n", e, metadata.FromContext(ctx))

	// 時間のかかる処理: 1 sec ~ 9.999 sec
	delay := 1000*time.Millisecond + time.Duration(rand.Intn(9000))*time.Millisecond
	begin := time.Now()
	for now := range time.Tick(1 * time.Second) {
		log.Printf(".")

		// time.Tickで取得した現在時間とループ開始直前の時間の差分でループを止めるか決める
		if now.Sub(begin) >= delay {
			break
		}
	}

	// 50% の確率で失敗
	if rand.Intn(100) < 50 {
		log.Println("fail")
		return fmt.Errorf("failed")
	}
	log.Println("pass")

	return nil
}

func main() {
	technical.Logger = logger.New(
		logger.WithOutput(os.Stdout),
		logger.WithTimeFormat(time.RFC3339Nano),
		logger.WithSkipFrameCount(3),
		logger.WithLevel(logger.InfoLevel),
	)

	b := rabbitmq.NewBroker(
		rabbitmq.ExchangeName(exchange),
		rabbitmq.PrefetchCount(1),
	)

	service := micro.NewService(
		micro.Broker(b),
		micro.WrapSubscriber(plogger.NewSubscriberWrapper()),
	)
	service.Init()

	subopts := broker.NewSubscribeOptions(
		rabbitmq.DurableQueue(), //
		broker.DisableAutoAck(),
		rabbitmq.RequeueOnError(),
		rabbitmq.AckOnSuccess(),
		rabbitmq.QueueArguments(map[string]interface{}{
			"x-max-priority": 10,
		}),
	)

	// register
	micro.RegisterSubscriber(topic, service.Server(), subscribe, server.SubscriberContext(subopts.Context), server.SubscriberQueue(queue))

	if err := service.Run(); err != nil {
		log.Fatalf("%+v", err)
	}
}
