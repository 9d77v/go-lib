package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/9d77v/go-lib/clients/config"
	"github.com/9d77v/go-lib/clients/etcd"
	"github.com/9d77v/go-lib/clients/jaeger"
	"github.com/opentracing-contrib/go-amqp/amqptracer"
	opentracing "github.com/opentracing/opentracing-go"
	tags "github.com/opentracing/opentracing-go/ext"
	"github.com/streadway/amqp"
)

//FailOnError mq fail
func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

//Client ..
type Client struct {
	Conn   *amqp.Connection
	Chs    map[string]*amqp.Channel
	Queues map[string]amqp.Queue
	Tracer opentracing.Tracer
	Closer io.Closer
}

//NewClient ..
func NewClient(config *config.RabbitmqConfig, queueNames []string) (*Client, error) {
	if len(queueNames) == 0 {
		log.Panicln("no queue")
	}
	client := new(Client)
	client.Chs = make(map[string]*amqp.Channel)
	client.Queues = make(map[string]amqp.Queue)
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/",
		config.User,
		config.Password,
		config.Host,
		config.Port))
	FailOnError(err, "Failed to connect to RabbitMQ")
	client.Conn = conn
	for _, v := range queueNames {
		ch, err := conn.Channel()
		FailOnError(err, "Failed to open a channel")
		client.Chs[v] = ch
		q, err := ch.QueueDeclare(
			v,     // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		FailOnError(err, "Failed to declare a queue")
		client.Queues[v] = q
	}
	return client, err
}

//NewClientFromEtcd init gorm from etcd config and watch config to update gorm
func NewClientFromEtcd(etcdCli *etcd.Client, queueNames []string) (mqCli *Client, err error) {
	appName := os.Getenv("APP_NAME")
	profile := os.Getenv("PROFILE")
	mqKey := etcdCli.GetEtcdKey(profile, appName, "rabbitmq")
	mqConfig := new(config.RabbitmqConfig)
	err = etcdCli.GetValue(5*time.Second, mqKey, mqConfig)
	if err != nil {
		log.Println("rabbitmq config is not exist:", err)
	}
	mq, err := NewClient(mqConfig, queueNames)
	if err != nil {
		log.Println("rabbitmq connect failed")
	}
	tracer, closer, err := jaeger.InitTracerFromEtcd(etcdCli, "rabbitmq")
	mq.Tracer = tracer
	mq.Closer = closer
	mqCli = mq
	log.Println("rabbitmq inited", mqCli)
	//change to new mq connection when  config changed
	go etcdCli.WatchKey(mqKey, mqConfig, func() {
		mq, err := NewClient(mqConfig, queueNames)
		if err != nil {
			log.Println("rabbitmq connect failed")
			return
		}
		for _, v := range mqCli.Chs {
			v.Close()
		}
		mq.Closer = mqCli.Closer
		mq.Tracer = mqCli.Tracer
		mqCli.Conn.Close()
		mqCli = mq
		log.Println("rabbitmq changed", mqCli)
	})
	return mqCli, err
}

//PublishMessage ..
func (c *Client) PublishMessage(ctx context.Context, exhcange, key string, immediate bool, msg amqp.Publishing) error {
	var span opentracing.Span
	if span = opentracing.SpanFromContext(ctx); span != nil {
		span = c.Tracer.StartSpan("PublishMessage", opentracing.ChildOf(span.Context()))
		tags.SpanKindProducer.Set(span)
		tags.PeerService.Set(span, "rabbitmq")
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	msg.Headers = make(map[string]interface{})
	if err := amqptracer.Inject(span, msg.Headers); err != nil {
		if span != nil {
			tags.Error.Set(span, true)
			span.LogKV("error", "amqptracer.Inject error")
		}
		return err
	}
	ch := c.Chs[key]
	if ch == nil {
		if span != nil {
			tags.Error.Set(span, true)
			span.LogKV("error", "channel is not exist")
		}
		return errors.New("channel is not exist")
	}
	return ch.Publish(exhcange, key, false, immediate, msg)
}
