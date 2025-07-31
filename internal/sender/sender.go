package sender

import (
	"encoding/json"
	"fmt"
	"github.com/berkinyildiran/insider-case/internal/cache"
	"github.com/berkinyildiran/insider-case/internal/message"
	"github.com/berkinyildiran/insider-case/internal/transporter"
	"log"
	"time"
)

type Sender struct {
	cache       cache.Cache
	config      *Config
	repository  *message.Repository
	transporter transporter.Transporter
}

func NewSender(cache cache.Cache, config *Config, repository *message.Repository, transporter transporter.Transporter) *Sender {
	return &Sender{
		cache:       cache,
		config:      config,
		repository:  repository,
		transporter: transporter,
	}
}

func (s *Sender) Run() error {
	messages, err := s.repository.GetPending(s.config.FetchLimit)
	if err != nil {
		return fmt.Errorf("failed to fetch pending messages (limit %d): %w", s.config.FetchLimit, err)
	}

	for _, msg := range messages {
		requestPayload, err := json.Marshal(RequestPayload{
			To:      msg.RecipientPhoneNumber,
			Content: msg.Content,
		})

		if err != nil {
			log.Printf("[ERROR] Failed to marshal request payload %s: %v", msg.ID, err)

			if err := s.repository.UpdateSendingStatus(msg.ID, message.FailedStatus); err != nil {
				log.Printf("[ERROR] Failed to update message %s as failed: %v", msg.ID, err)
			}

			continue
		}

		response, err := s.transporter.Send(s.config.Address, requestPayload)
		sendingTime := time.Now()

		if err != nil {
			log.Printf("[ERROR] Failed to send message %s: %v", msg.ID, err)

			if err := s.repository.UpdateSendingStatus(msg.ID, message.FailedStatus); err != nil {
				log.Printf("[ERROR] Failed to update message %s as failed: %v", msg.ID, err)
			}

			continue
		}

		if err := s.repository.UpdateSendingStatus(msg.ID, message.SuccessStatus); err != nil {
			log.Printf("[ERROR] Failed to update message %s as success: %v", msg.ID, err)
			continue
		}

		log.Printf("[INFO] Successfully sent message %s", msg.ID)

		var responsePayload ResponsePayload
		if err := json.Unmarshal(response, &responsePayload); err != nil {
			log.Printf("[ERROR] Failed to unmarshal response payload %s: %v", msg.ID, err)
			continue
		}

		cachePayload, err := json.Marshal(CachePayload{
			SendingTime: sendingTime,
			MessageID:   responsePayload.MessageID,
		})

		if err = s.cache.Set(msg.ID.String(), cachePayload, time.Duration(s.config.CacheTTL)*time.Second); err != nil {
			log.Printf("[ERROR] Failed to cache response for message %s: %v", msg.ID, err)
		}

		log.Printf("[INFO] Successfully cached response for message %s", msg.ID)
	}

	return nil
}
