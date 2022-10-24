package rabbitmq

import (
	"encoding/json"
)

type messageType uint64

const (
	messageSend messageType = iota + 1
	messageClose
)

type message struct {
	CorrelationID string
	Type          messageType
	Data          json.RawMessage
}
