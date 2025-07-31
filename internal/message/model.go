package message

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Message struct {
	ID                   uuid.UUID     `gorm:"type:uuid"`
	Content              string        `gorm:"type:varchar(100);not null"`
	RecipientPhoneNumber string        `gorm:"type:varchar(20);not null"`
	SendingStatus        SendingStatus `gorm:"type:smallint;not null;index:idx_messages_sending_status_created_at,priority:1"`
	CreatedAt            time.Time     `gorm:"index:idx_messages_sending_status_created_at,priority:2,sort:asc"`
	UpdatedAt            time.Time
}

func (m *Message) BeforeCreate(tx *gorm.DB) error {
	m.ID = uuid.New()
	m.SendingStatus = PendingStatus

	return nil
}
