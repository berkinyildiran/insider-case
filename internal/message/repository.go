package message

import (
	"context"
	"github.com/berkinyildiran/insider-case/internal/database"
	"github.com/google/uuid"
)

type Repository struct {
	database *database.Database

	context context.Context
}

func NewRepository(database *database.Database, context context.Context) *Repository {
	return &Repository{
		database: database,

		context: context,
	}
}

func (r *Repository) GetPending(limit int) ([]Message, error) {
	r.database.WG.Add(1)
	defer r.database.WG.Done()

	var messages []Message
	result := r.database.DB.
		WithContext(r.context).
		Where("sending_status", PendingStatus).
		Order("created_at asc").
		Limit(limit).
		Find(&messages)

	return messages, result.Error
}

func (r *Repository) GetSent(limit int, offset int) ([]Message, error) {
	r.database.WG.Add(1)
	defer r.database.WG.Done()

	var messages []Message
	result := r.database.DB.
		WithContext(r.context).
		Where("sending_status", SuccessStatus).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&messages)

	return messages, result.Error
}

func (r *Repository) UpdateSendingStatus(id uuid.UUID, status SendingStatus) error {
	r.database.WG.Add(1)
	defer r.database.WG.Done()

	result := r.database.DB.
		WithContext(r.context).
		Model(&Message{}).
		Where("id", id).
		Update("sending_status", status)

	return result.Error
}
