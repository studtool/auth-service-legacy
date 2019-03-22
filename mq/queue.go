package mq

import (
	"auth-service/config"
	"auth-service/errs"
	"fmt"
	"github.com/streadway/amqp"
)

type MQ struct {
	cq      amqp.Queue
	ch      *amqp.Channel
	conn    *amqp.Connection
	connStr string
}

func NewQueue() *MQ {
	return &MQ{
		connStr: fmt.Sprintf("amqp://%s:%s@%s:%s/",
			config.MessageQueueUser, config.MessageQueuePassword,
			config.MessageQueueHost, config.MessageQueuePort,
		),
	}
}

func (mq *MQ) OpenConnection() error {
	conn, err := amqp.Dial(mq.connStr)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		config.CreatedUsersQueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	mq.cq = q
	mq.ch = ch
	mq.conn = conn

	return nil
}

func (mq *MQ) CloseConnection() {
	_ = mq.ch.Close()
	_ = mq.conn.Close()
}

func (mq *MQ) SendUserCreated(userId string) *errs.Error {
	err := mq.ch.Publish(
		"",
		mq.cq.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(userId),
		},
	)
	if err != nil {
		return errs.NewInternalError(err.Error())
	}
	return nil
}
