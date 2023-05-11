package replyButton

import (
	"multiBot/config"
	"multiBot/errorCodes"
	"multiBot/models"
	makeReplyButtonMessage "multiBot/sender/replyButton"
	"multiBot/types/replyButton"
	"multiBot/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

func ValiadateReplyButtonPayload(c *fiber.Ctx) error {
	var replyButtonPayload []replyButton.ReplyButtonPayloadClient
	err := c.BodyParser(&replyButtonPayload)
	if err != nil {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: err.Error(),
				ErrorCode:    errorCodes.REQUEST_BODY_PARSING_ERROR,
			},
			replyButtonPayload,
		)
	}

	for _, msg := range replyButtonPayload {
		v := validate.Struct(msg)
		if !v.Validate() {
			return utils.ErrorResponse(
				c,
				utils.AppError{
					DebugMessage: v.Errors.One(),
					ErrorCode:    errorCodes.VALIDATION_ERROR,
				},
				replyButtonPayload,
			)
		}
	}

	cfg := config.GetConfig()
	messageLimit := cfg.API_MESSAGE_LIMIT
	if replyButtonPayload != nil && len(replyButtonPayload) > messageLimit {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: "",
				ErrorCode:    errorCodes.MESSAGE_LIMIT_EXCEEDED,
			},
			replyButtonPayload,
		)
	}

	c.Locals("replyButtonPayload", replyButtonPayload)
	return c.Next()
}

func ReplyButton(c *fiber.Ctx) error {
	replyButtonPayload := c.Locals("replyButtonPayload").([]replyButton.ReplyButtonPayloadClient)

	var msgSlice []models.Message
	clientIdentifier := c.Locals("clientIdentifier").(string)
	waAccountIdentifier := c.Locals("waAccountIdentifier").(string)

	for _, msg := range replyButtonPayload {
		replyButtonMessage, err := makeReplyButtonMessage.MakeReplyButtonMessage(
			msg.Body,
			msg.To,
			msg.Purpose,
			clientIdentifier,
			waAccountIdentifier,
		)
		if err != nil {
			return utils.ErrorResponse(
				c,
				utils.AppError{
					DebugMessage: err.Error(),
					ErrorCode:    errorCodes.MAKE_MESSAGE_ERROR,
				},
				replyButtonPayload,
			)
		}
		msgSlice = append(msgSlice, *replyButtonMessage)
	}

	var messageModel models.Message
	objIds, err := messageModel.InsertManyMessages(msgSlice)
	if err != nil {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: err.Error(),
				ErrorCode:    errorCodes.DB_OPERATION_ERROR,
			},
			replyButtonPayload,
		)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"ids": objIds,
	}, replyButtonPayload)
}
