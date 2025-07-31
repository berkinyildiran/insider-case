package sender

import (
	"github.com/google/uuid"
	"time"
)

type RequestPayload struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

type ResponsePayload struct {
	Message   string    `json:"message"`
	MessageID uuid.UUID `json:"messageId"`
}

type CachePayload struct {
	SendingTime time.Time `json:"sendingTime"`
	MessageID   uuid.UUID `json:"messageId"`
}
