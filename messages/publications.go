package messages

import (
	"github.com/studtool/common/errs"
	"github.com/studtool/common/queues"
)

func (c *MqClient) PublishUserCreatedMessage(
	data *queues.CreatedUserData,
) *errs.Error {
	err := c.publishMessage(queues.CreatedUsersQueueName, data)
	if err != nil {
		return errs.New(err)
	}
	return nil
}

func (c *MqClient) PublishUserDeletedMessage(
	data *queues.DeletedUserData,
) *errs.Error {
	err := c.publishMessage(queues.DeletedUsersQueueName, data)
	if err != nil {
		return errs.New(err)
	}
	return nil
}

func (c *MqClient) PublishRegistrationEmailMessage(
	data *queues.RegistrationEmailData,
) *errs.Error {
	err := c.publishMessage(queues.RegistrationEmailsQueueName, data)
	if err != nil {
		return errs.New(err)
	}
	return nil
}
