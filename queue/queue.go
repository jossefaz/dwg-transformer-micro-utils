package queue

import (
	"fmt"
	"github.com/streadway/amqp"
	"os"
	"github.com/yossefazoulay/go_utils/utils"
)

type rabbitmq struct {
	Conn  *amqp.Connection
	ChanL *amqp.Channel
	Queue amqp.Queue
}


func NewRabbit(connString string, queueName string) (instance rabbitmq) {
	conn := dial(connString)
	amqpChannel := getChannel(conn)
	queue := connectToQueue(amqpChannel, queueName)
	return rabbitmq{
		Conn:  conn,
		ChanL: amqpChannel,
		Queue: queue,
	}
}

func dial(connString string) *amqp.Connection {
	conn, err := amqp.Dial(connString)
	utils.HandleError(err, "Can't connect to AMQP")
	if err != nil {
		os.Exit(1)
	}
	return conn
}

func getChannel(conn *amqp.Connection) *amqp.Channel {
	c, err := conn.Channel()
	utils.HandleError(err, "Can't create a amqpChannel")
	if err != nil {
		os.Exit(1)
	}
	return c
}

func connectToQueue(c *amqp.Channel, queueName string) amqp.Queue {
	q, err := c.QueueDeclare(queueName, true, false, false, false, nil)
	utils.HandleError(err, "Could not declare `add` queue")
	if err != nil {
		os.Exit(1)
	}
	return q
}

func (rmq rabbitmq) SendMessage(body []byte) {
	err := rmq.ChanL.Publish("", rmq.Queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}
func (rmq rabbitmq) ListenMessage(onMessage func(m amqp.Delivery, rmq rabbitmq)) {
	err := rmq.ChanL.Qos(1, 0, false)
	utils.HandleError(err, "Could not configure QoS")
	messageChannel, err := rmq.ChanL.Consume(
		rmq.Queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	utils.HandleError(err, "Could not register consumer")
	stopChan := make(chan bool)
	go func() {
		for d := range messageChannel {
			onMessage(d, rmq)
		}
	}()
	// Stop for program termination
	<-stopChan

}
