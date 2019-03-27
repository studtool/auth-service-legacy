package mq

import (
	"auth-service/config"
	"auth-service/errs"
	"fmt"
	"github.com/streadway/amqp"
)

type MQ struct {
	cq      amqp.Queue
	dq      amqp.Queue
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

	mq.cq, err = ch.QueueDeclare(
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

	mq.dq, err = ch.QueueDeclare(
		config.DeletedUsersQueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	mq.ch = ch
	mq.conn = conn

	return nil
}

func (mq *MQ) CloseConnection() error {
	if err := mq.ch.Close(); err != nil {
		return err
	}
	return mq.conn.Close()
}

func (mq *MQ) SendUserCreated(userId string) *errs.Error {
	return mq.sendUserId(mq.cq, userId)
}

func (mq *MQ) SendUserDeleted(userId string) *errs.Error {
	return mq.sendUserId(mq.dq, userId)
}

func (mq *MQ) sendUserId(q amqp.Queue, userId string) *errs.Error {
	err := mq.ch.Publish(
		"",
		q.Name,
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