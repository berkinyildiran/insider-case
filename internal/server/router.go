package server

import (
	"context"
	"fmt"
	"github.com/berkinyildiran/insider-case/internal/message"
	"github.com/berkinyildiran/insider-case/internal/validator"
	"github.com/gofiber/fiber/v3"
)

type Router struct {
	config    *Config
	handler   *message.Handler
	validator *validator.Validator

	app     *fiber.App
	context context.Context
}

func NewRouter(config *Config, handler *message.Handler, validator *validator.Validator, context context.Context) *Router {
	return &Router{
		config:    config,
		handler:   handler,
		validator: validator,

		app:     fiber.New(),
		context: context,
	}
}

func (r *Router) Setup() {
	messaging := r.app.Group("/messaging")

	messaging.Get("/sent", r.handler.GetSentMessages)
	messaging.Post("/start", r.handler.StartScheduler)
	messaging.Post("/stop", r.handler.StopScheduler)
}

func (r *Router) Start() error {
	address := fmt.Sprintf(":%d", r.config.Port)
	return r.app.Listen(address)
}

func (r *Router) Stop() error {
	return r.app.Shutdown()
}
