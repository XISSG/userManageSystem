package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"log"
)

func readConfig(filename string) *Config {
	viper.AddConfigPath("./conf")
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var config *Config
	err = viper.Unmarshal(&config)
	return config
}

func initRabbitMQ() (*amqp.Connection, error) {
	config := readConfig("rabbitmq")
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%v/", config.User, config.Password, config.Host, config.Port)
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

type RabbitMQService struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue

	chg Exchange
	key string
}

func NewRabbitQMService() *RabbitMQService {
	conn, err := initRabbitMQ()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &RabbitMQService{
		conn: conn,
	}
}
func (r *RabbitMQService) DeclareChannel() error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	r.ch = ch
	return nil
}

func (r *RabbitMQService) DeclareExChange(chg Exchange) error {
	// 声明一个交换机
	err := r.ch.ExchangeDeclare(
		chg.name,       // 交换机名称
		chg.kind,       // 交换机类型
		chg.durable,    // 是否持久化
		chg.autoDelete, // 是否自动删除
		chg.internal,   // 是否内部使用
		chg.noWait,     // 是否等待服务器响应
		chg.args,       // 其他属性
	)
	r.chg = chg
	return err
}

func (r *RabbitMQService) DeclareQueue(que Queue) error {
	q, err := r.ch.QueueDeclare(
		que.name,       // 队列名称
		que.durable,    // 是否持久化
		que.autoDelete, // 是否自动删除
		que.exclusive,  // 是否排他
		que.noWait,     // 是否等待服务器响应
		que.args,       // 其他属性
	)
	r.q = q
	return err
}

func (r *RabbitMQService) BindQueue(bind Bind) error {
	err := r.ch.QueueBind(
		bind.name,
		bind.key,
		bind.exchange, // 交换机名称
		bind.noWait,
		bind.args, // 其他属性
	)
	return err
}
func (r *RabbitMQService) Close() error {
	return r.conn.Close()
}

func (r *RabbitMQService) Publish(publish Publish) error {
	err := r.ch.Publish(
		publish.exchange,
		publish.key,
		publish.immediate,
		publish.immediate,
		publish.msg,
	)
	return err
}
func (r *RabbitMQService) Consume(consume Consume) ([]string, error) {
	msgs, err := r.ch.Consume(
		consume.queue,
		consume.consumer,
		consume.autoAck,
		consume.exclusive,
		consume.noLocal,
		consume.noWait,
		consume.args, // 其他属性
	)
	if err != nil {
		return nil, err
	}

	var res []string
	for d := range msgs {
		res = append(res, string(d.Body))
	}
	return res, nil
}

func NewDefaultMQPublisher(data string) {
	service := NewRabbitQMService()
	err := service.DeclareChannel()
	if err != nil {
		log.Fatal(err)
		return
	}
	que := Queue{
		"default", // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,
	}
	err = service.DeclareQueue(que)
	if err != nil {
		log.Fatal(err)
		return
	}
	publish := Publish{
		"",             // exchange
		service.q.Name, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		},
	}
	err = service.Publish(publish)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func NewDefaultMQConsumer() []string {
	service := NewRabbitQMService()
	err := service.DeclareChannel()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	que := Queue{
		"default", // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,
	}
	err = service.DeclareQueue(que)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	consume := Consume{
		service.q.Name, // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	}
	msgs, err := service.Consume(consume)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return msgs
}
