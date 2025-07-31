package http

import (
	"fmt"
	"github.com/berkinyildiran/insider-case/internal/transporter"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type Http struct {
	client *resty.Client
}

func NewHttp() *Http {
	return &Http{
		client: resty.New(),
	}
}

func (h *Http) Send(address string, payload transporter.Payload) (transporter.Response, error) {
	response, err := h.client.
		NewRequest().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post(address)

	if err != nil {
		return nil, fmt.Errorf("failed to execute POST request to %s: %w", address, err)
	}

	if response.StatusCode() != http.StatusOK && response.StatusCode() != http.StatusCreated && response.StatusCode() != http.StatusAccepted {
		return nil, fmt.Errorf("unexpected status code %d (%s) from %s", response.StatusCode(), response.Status(), address)
	}

	return response.Body(), nil
}
