package message

import (
	"github.com/berkinyildiran/insider-case/internal/scheduler"
	"github.com/berkinyildiran/insider-case/internal/validator"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	repository *Repository
	scheduler  *scheduler.Scheduler
	validator  *validator.Validator
}

func NewHandler(repository *Repository, scheduler *scheduler.Scheduler, validator *validator.Validator) *Handler {
	return &Handler{
		repository: repository,
		scheduler:  scheduler,
		validator:  validator,
	}
}

// StartScheduler godoc
// @Summary      Start scheduler
// @Description  Starts the scheduler that sends messages periodically
// @Tags         messaging
// @Success      200  {object}  GenericSuccessResponse
// @Router       /messaging/start [post]
func (h *Handler) StartScheduler(c fiber.Ctx) error {
	message := h.scheduler.Start()

	return c.JSON(GenericSuccessResponse{
		Message: message,
	})
}

// StopScheduler godoc
// @Summary      Stop message scheduler
// @Description  Stops the running message scheduler
// @Tags         messaging
// @Success      200  {object}  GenericSuccessResponse
// @Router       /messaging/stop [post]
func (h *Handler) StopScheduler(c fiber.Ctx) error {
	message := h.scheduler.Stop()

	return c.JSON(GenericSuccessResponse{
		Message: message,
	})
}

// GetSentMessages godoc
// @Summary      Get sent messages
// @Description  Retrieves sent messages with pagination support
// @Tags         messaging
// @Accept       json
// @Param        limit   query     int    false  "Limit number of messages" minimum(1) maximum(100)
// @Param        offset  query     int    false  "Pagination offset" minimum(0)
// @Success      200     {array}   message.Message
// @Failure      400     {object}  GenericFailureResponse
// @Failure      500     {object}  GenericFailureResponse
// @Router       /messaging/sent [get]
func (h *Handler) GetSentMessages(c fiber.Ctx) error {
	var query GetSendMessagesQuery
	if err := c.Bind().Query(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(GenericFailureResponse{
			Success: false,
			Error:   "invalid query parameters",
		})
	}

	if err := h.validator.Validate(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(GenericFailureResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	messages, err := h.repository.GetSent(query.Limit, query.Offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(GenericFailureResponse{
			Success: false,
			Error:   "failed to fetch sent messages",
		})
	}

	return c.JSON(messages)
}
