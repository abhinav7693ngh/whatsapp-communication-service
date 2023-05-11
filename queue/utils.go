package queue

import (
	kafkaGo "github.com/segmentio/kafka-go"
)

var MessageChannel chan kafkaGo.Message = make(chan kafkaGo.Message)
