package queue

import (
	"fmt"
	"github.com/streadway/amqp"
	"os"
	"github.com/yossefazoulay/go_utils/utils"
)



type queues map[string]amqp.Queue

type Rabbitmq struct {
	Conn  *amqp.Connection
	ChanL *amqp.Channel
	Queues queues

}


func NewRabbit(connString string, queuesName []string) (Rabbitmq, error) {
	conn, err := dial(connString)
	if err != nil {
		return Rabbitmq{}, err
	}
	amqpChannel, err := getChannel(conn)
	if err != nil {
		return Rabbitmq{}, err
	}
	qs, err := declareQueues(amqpChannel, queuesName)
	if err != nil {
		return Rabbitmq{}, err
	}
	return Rabbitmq{
		Conn:  conn,
		ChanL: amqpChannel,
		Queues: qs,
	}, nil
}

func declareQueues(c *amqp.Channel, queuesName []string) (queues, error){
	qs := queues{}
	for _, qu := range queuesName {
		q, err := connectToQueue(c, qu)
		if err != nil {
			return qs, err
		}
		qs[qu] = q
	}
	return qs, nil
}

func dial(connString string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(connString)
	if err != nil {
		return conn, err
	}
	return conn, nil
}

func getChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	c, err := conn.Channel()
	if err != nil {
		return c, err
	}
	return c, nil
}

func connectToQueue(c *amqp.Channel, queueName string) (amqp.Queue, error) {
	q, err := c.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return q, err
	}
	return q, nil
}

func (rmq Rabbitmq) SendMessage(body []byte, queueName string) {
	err := rmq.ChanL.Publish("", rmq.Queues[queueName].Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}
func (rmq Rabbitmq) ListenMessage(onMessage func(m amqp.Delivery, q Rabbitmq), queueName string) error {
	err := rmq.ChanL.Qos(1, 0, false)
	if err != nil {
		return err
	}
	messageChannel, err := rmq.ChanL.Consume(
		rmq.Queues[queueName].Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	stopChan := make(chan bool)
	go func() {
		for d := range messageChannel {
			onMessage(d, rmq)
		}
	}()

	// Stop for program termination
	<-stopChan
	return nil

}

func (rmq Rabbitmq) OpenListening (c []string, cb func(m amqp.Delivery, q Rabbitmq)) {
	for _, q := range c {
		rmq.ListenMessage(cb, q)
	}
}