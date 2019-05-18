package messages

import (
	"fmt"

	"github.com/streadway/amqp"

	"github.com/studtool/common/queues"
	"github.com/studtool/common/utils"

	"github.com/studtool/auth-service/beans"
	"github.com/studtool/auth-service/config"
)

type MqClient struct {
	connStr    string
	connection *amqp.Connection

	channel *amqp.Channel

	registrationEmailsQueue amqp.Queue
	createdUsersQueue       amqp.Queue
	deletedUsersQueue       amqp.Queue
}

func NewMqClient() *MqClient {
	return &MqClient{
		connStr: fmt.Sprintf("amqp://%s:%s@%s:%d/",
			config.MqUser.Value(), config.MqPassword.Value(),
			config.MqHost.Value(), config.MqPort.Value(),
		),
	}
}

func (c *MqClient) OpenConnection() error {
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

	c.connection = conn

	c.channel, err = conn.Channel()
	if err != nil {
		return err
	}

	c.createdUsersQueue, err =
		c.declareQueue(queues.CreatedUsersQueueName)
	if err != nil {
		return err
	}

	c.deletedUsersQueue, err =
		c.declareQueue(queues.DeletedUsersQueueName)
	if err != nil {
		return err
	}

	c.registrationEmailsQueue, err =
		c.declareQueue(queues.RegistrationEmailsQueueName)
	if err != nil {
		return err
	}

	return nil
}

func (c *MqClient) CloseConnection() error {
	if err := c.channel.Close(); err != nil {
		return err
	}
	return c.connection.Close()
}

func (c *MqClient) writeEnqueueLog(queue string, data interface{}) {
	beans.Logger().Info(fmt.Sprintf("queue: %s <- %v", queue, data))
}
