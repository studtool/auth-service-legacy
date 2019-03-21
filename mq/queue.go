package mq

import (
	"github.com/streadway/amqp"
)

type MQ struct {
}

func NewQueue() *MQ {
	return &MQ{}
}

func (q *MQ) Open() {
	_, _ = amqp.Dial("amqp://guest:guest@localhost:5672/") //TODO
}

func (q *MQ) Close() {
	//TODO
}
