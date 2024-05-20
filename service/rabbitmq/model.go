package rabbitmq

import "github.com/rabbitmq/amqp091-go"

type Config struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

type Context struct {
	exchage Exchange
	que     Queue
}

type Exchange struct {
	name       string        // 交换机名称
	kind       string        // 交换机类型
	durable    bool          // 是否持久化
	autoDelete bool          // 是否自动删除
	internal   bool          // 是否内部使用
	noWait     bool          // 是否等待服务器响应
	args       amqp091.Table // 其他属性
}

type Queue struct {
	name       string
	durable    bool
	autoDelete bool
	exclusive  bool
	noWait     bool
	args       amqp091.Table // 其他属性
}

type Bind struct {
	name     string
	key      string
	exchange string
	noWait   bool
	args     amqp091.Table // 其他属性
}

type Publish struct {
	exchange  string
	key       string
	mandatory bool
	immediate bool
	msg       amqp091.Publishing
}

type Consume struct {
	queue     string
	consumer  string
	autoAck   bool
	exclusive bool
	noLocal   bool
	noWait    bool
	args      amqp091.Table
}
