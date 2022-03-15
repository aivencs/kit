package main

import (
	"context"
	"fmt"

	"github.com/aivencs/kit/pkg/messenger"
	"github.com/streadway/amqp"
)

func main() {
	ctx := context.WithValue(context.Background(), "trace", "ctx-messenger-001")
	option := messenger.MenssengerOption{
		Host:      "localhost:5672",
		Auth:      false,
		Username:  "",
		Password:  "",
		Zone:      "/",
		Topics:    map[string]string{"consume": "article_draft", "product": "article_draft_eh", "forward": "", "bad": ""},
		Heartbeat: 120,
		Qos:       1,
	}
	messenger.InitMessenger("rabbit", option)
	if !messenger.GetActive(ctx) {
		fmt.Println("not active")
		return
	}
	conn := messenger.GetConnect(ctx).(*amqp.Connection)
	chl, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	consume, err := messenger.GetConsume(ctx, chl)
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case carton := <-consume.(<-chan amqp.Delivery):
			fmt.Println(map[string]interface{}{"p": carton.Priority, "m": string(carton.Body)})
			carton.Ack(false)
			messenger.Sent(ctx, messenger.SentPayload{
				Topic:    messenger.GetTopic(ctx, "product"),
				Message:  fmt.Sprintf("abc-%s", string(carton.Body)),
				Priority: 5,
				Channel:  chl,
			})
		default:
			continue
		}
	}
}
