package messenger

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

var once sync.Once
var messenger Messenger

type Messenger interface {
	Sent(ctx context.Context, payload SentPayload) error
	GetConnect(ctx context.Context) interface{}
	GetTopic(ctx context.Context, key string) string
	GetConsume(ctx context.Context, channel interface{}) (interface{}, error)
	GetActive(ctx context.Context) bool
}

type SentPayload struct {
	Topic    string
	Message  string
	Priority uint8
	Channel  interface{}
}

type Rabbit struct {
	Connect *amqp.Connection
	Topics  map[string]string
	rwMutex *sync.RWMutex
	Active  bool
	Qos     int
	Channel *amqp.Channel
}

type MenssengerOption struct {
	Host      string
	Auth      bool
	Zone      string
	Username  string
	Password  string
	Topics    map[string]string
	Heartbeat int
	Active    bool
	Qos       int
}

func InitMessenger(name string, opt MenssengerOption) {
	once.Do(func() {
		switch name {
		case "rabbit":
			messenger = NewRabbit(opt)
		default:
			messenger = NewRabbit(opt)
		}
	})
}

func NewRabbit(opt MenssengerOption) Messenger {
	conf := amqp.Config{
		Heartbeat: time.Second * time.Duration(opt.Heartbeat),
	}
	address := fmt.Sprintf("amqp://%s%s", opt.Host, opt.Zone)
	if opt.Auth {
		address = fmt.Sprintf("amqp://%s:%s@%s%s", opt.Username, opt.Password, opt.Host, opt.Zone)
	}
	conn, err := amqp.DialConfig(address, conf)
	if err != nil {
		return &Rabbit{Active: false, rwMutex: new(sync.RWMutex)}
	}
	return &Rabbit{
		Active:  true,
		Topics:  opt.Topics,
		Connect: conn,
		rwMutex: new(sync.RWMutex),
		Qos:     opt.Qos,
	}
}

func (c *Rabbit) Sent(ctx context.Context, payload SentPayload) error {
	chl := payload.Channel.(*amqp.Channel)
	err := chl.Publish(
		c.Topics["product"], "", false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(payload.Message),
			Priority:    payload.Priority,
		},
	)
	return err
}

func (c *Rabbit) GetConnect(ctx context.Context) interface{} {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()
	return c.Connect
}

func (c *Rabbit) GetActive(ctx context.Context) bool {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()
	return c.Active
}

func (c *Rabbit) GetTopic(ctx context.Context, key string) string {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()
	return c.Topics[key]
}

func (c *Rabbit) GetConsume(ctx context.Context, channel interface{}) (interface{}, error) {
	conn := c.GetConnect(ctx).(*amqp.Connection)
	chl, err := conn.Channel()
	if err != nil {
		return chl, err
	}
	err = chl.Qos(c.Qos, 0, false)
	if err != nil {
		return chl, err
	}
	return chl.Consume(
		c.Topics["consume"], // queue
		"",                  // consumer
		false,               // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
}

func GetConsume(ctx context.Context, channel interface{}) (interface{}, error) {
	return messenger.GetConsume(ctx, channel)
}

func GetTopic(ctx context.Context, key string) string {
	return messenger.GetTopic(ctx, key)
}

func GetConnect(ctx context.Context) interface{} {
	return messenger.GetConnect(ctx)
}

func Sent(ctx context.Context, payload SentPayload) error {
	return messenger.Sent(ctx, payload)
}

func GetActive(ctx context.Context) bool {
	return messenger.GetActive(ctx)
}
