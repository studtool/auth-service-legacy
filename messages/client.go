package messages

import (
	"fmt"

	"github.com/mailru/easyjson"
	"github.com/streadway/amqp"

	"github.com/studtool/common/consts"
	"github.com/studtool/common/errs"
	"github.com/studtool/common/queues"
	"github.com/studtool/common/utils"

	"github.com/studtool/auth-service/beans"
	"github.com/studtool/auth-service/config"
)

type QueueClient struct {
	connStr    string
	connection *amqp.Connection

	channel *amqp.Channel

	regEmailsQueue    amqp.Queue
	createdUsersQueue amqp.Queue
	deletedUsersQueue amqp.Queue
}

func NewQueueClient() *QueueClient {
	return &QueueClient{
		connStr: fmt.Sprintf("amqp://%s:%s@%s:%d/",
			config.MqUser.Value(), config.MqPassword.Value(),
			config.MqHost.Value(), config.MqPort.Value(),
		),
	}
}

func (c *QueueClient) OpenConnection() error {
	var conn *amqp.Connection
	err := utils.WithRetry(func(n int) (err error) {
		if n > 0 {
			beans.Logger().Info(fmt.Sprintf("opening message queue connection. retry #%d", n))
		}
		conn, err = amqp.Dial(c.connStr)
		return err
	}, config.MqConnNumRet.Value(), config.MqConnRetItv.Value())
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	c.regEmailsQueue, err = ch.QueueDeclare(
		queues.RegistrationEmailsQueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	c.createdUsersQueue, err = ch.QueueDeclare(
		queues.CreatedUsersQueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	c.channel = ch
	c.connection = conn

	return nil
}

func (c *QueueClient) CloseConnection() error {
	if err := c.channel.Close(); err != nil {
		return err
	}
	return c.connection.Close()
}

func (c *QueueClient) SendRegEmailMessage(data *queues.RegistrationEmailData) *errs.Error {
	body, _ := easyjson.Marshal(data)

	err := c.channel.Publish(
		consts.EmptyString,
		c.regEmailsQueue.Name,
		false,
		false,
		amqp.Publishing{
			Body:        body,
			ContentType: "application/json",
		},
	)
	if err != nil {
		return errs.New(err)
	}

	return nil
}

func (c *QueueClient) SendUserCreatedMessage(data *queues.CreatedUserData) *errs.Error {
	body, _ := easyjson.Marshal(data)

	err := c.channel.Publish(
		consts.EmptyString,
		c.createdUsersQueue.Name,
		false,
		false,
		amqp.Publishing{
			Body:        body,
			ContentType: "application/json",
		},
	)
	if err != nil {
		return errs.New(err)
	}

	return nil
}
