package queue

import (
	"fmt"
	"github.com/streadway/amqp"
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
	amqpChannel, err1 := getChannel(conn)
	if err1 != nil {
		return Rabbitmq{}, err1
	}
	qs, err2 := declareQueues(amqpChannel, queuesName)
	if err2 != nil {
		return Rabbitmq{}, err2
	}
	err3 := sendOnlyIfAck(amqpChannel)
	if err3 != nil {
		return Rabbitmq{}, err3
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

func sendOnlyIfAck(ch *amqp.Channel) error {
	err := ch.Qos(1, 0, false)
	if err != nil {
		return err
	}
}

func (rmq Rabbitmq) SendMessage(body []byte, queueName string, from string) (string, error) {
	err := rmq.ChanL.Publish("", rmq.Queues[queueName].Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
		Headers: map[string]interface{}{
			"From" : from,
			"To" : queueName,
		},
	})
	if err != nil {
		return "error occured when sending message : ", err

	}
	return string(body), nil

}
func (rmq *Rabbitmq) ListenMessage(onMessage func(m amqp.Delivery, q *Rabbitmq, queueName string), queueName string) error {

	rmq.ChanL, _ = getChannel(rmq.Conn)

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

	go func() {
		for d := range messageChannel {
			onMessage(d, rmq, queueName )
		}
	}()

	// Stop for program termination

	return nil

}

func (rmq *Rabbitmq) OpenListening (c []string, cb func(m amqp.Delivery, q *Rabbitmq, queueName string)) error {
	stopChan := make(chan bool)
	for _, q := range c {
		go func() {
			err := rmq.ListenMessage(cb, q)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
		}()
	}
	<-stopChan
	return nil
}