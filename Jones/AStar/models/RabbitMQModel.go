package models

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/astaxie/beego"
)

func failOnError(err error, msg string) {
	if err != nil {
		Logger.Error(msg + "\t" + err.Error())
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

type RabbitMQModel struct {
	channel *amqp.Channel
	queue   amqp.Queue
}

var mqInstance *RabbitMQModel

func (model RabbitMQModel) GetRabbitMQ() *RabbitMQModel {
	if mqInstance == nil {
		Logger.Info("Init RabbitMQ")
		urlMQ := beego.AppConfig.String("rabbitMQ")
		conn, err := amqp.Dial(urlMQ)
		failOnError(err, "Failed to connect to RabbitMQ")
		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		queueName := beego.AppConfig.String("MQName")
		// add exchange declare
		if exchangeErr := ch.ExchangeDeclare(queueName, "direct", true, false, false, false, nil);
			exchangeErr != nil {
			fmt.Println("exchange error:", exchangeErr)
		}
		q, err := ch.QueueDeclare(
			queueName, // name
			false,     // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		failOnError(err, "Failed to declare a queue")
		mqInstance = &RabbitMQModel{ch, q}
		return mqInstance
	}
	return mqInstance
}

func (model RabbitMQModel) Publish(data []byte) {
	/*var mqModel = model.GetRabbitMQ()
	err := mqInstance.channel.Publish(
		"",
		mqModel.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		})
	failOnError(err, "Failed to publish a message")*/
	model.PublishRetry(data, 2)
}

func (model RabbitMQModel) PublishRetry(data []byte, retry int) {
	var mqModel = model.GetRabbitMQ()
	err := mqInstance.channel.Publish(
		"",
		mqModel.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		})
	if err != nil {
		fmt.Println(err.Error())
		retry = retry - 1
		if retry <= 0 {
			failOnError(err, "Failed to publish a message")
		} else {
			fmt.Printf("mqInstance 置空, 重试: %d\n", retry)
			mqInstance = nil // 重新创建mq实例
			model.PublishRetry(data, retry)
		}
	}

}
