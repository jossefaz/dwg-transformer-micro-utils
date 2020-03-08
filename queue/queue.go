package queue

import (
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
func (rmq Rabbitmq) ListenMessage(onMessage func(m amqp.Delivery, q Rabbitmq), queueName string) (<-chan amqp.Delivery, error) {
	err := rmq.ChanL.Qos(1, 0, false)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	// Stop for program termination
	return messageChannel, nil

}

func (rmq Rabbitmq) OpenListening (c []string, cb func(m amqp.Delivery, q Rabbitmq)) error {
	stopChan := make(chan bool)
	sChannel := []<-chan amqp.Delivery{}
	for _, q := range c {
		mChan, err := rmq.ListenMessage(cb, q)
		if err != nil {
			return err
		}
		sChannel = append(sChannel, mChan)
	}

	go func() {
		for _, mch := range sChannel {
			for d := range mch {
				cb(d, rmq)
			}
		}
	}()
	<-stopChan
	return nil
}