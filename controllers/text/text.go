package text

import (
	"multiBot/config"
	"multiBot/errorCodes"
	"multiBot/models"
	makeTextMessage "multiBot/sender/text"
	"multiBot/types/text"
	"multiBot/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

func ValidateTextPayload(c *fiber.Ctx) error {
	var textPayload []text.TextPayloadClient
	err := c.BodyParser(&textPayload)
	if err != nil {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: err.Error(),
				ErrorCode:    errorCodes.REQUEST_BODY_PARSING_ERROR,
			},
			textPayload,
		)
	}

	for _, msg := range textPayload {
		v := validate.Struct(msg)
		if !v.Validate() {
			return utils.ErrorResponse(
				c,
				utils.AppError{
					DebugMessage: v.Errors.One(),
					ErrorCode:    errorCodes.VALIDATION_ERROR,
				},
				textPayload,
			)
		}
	}

	cfg := config.GetConfig()
	messageLimit := cfg.API_MESSAGE_LIMIT
	if textPayload != nil && len(textPayload) > messageLimit {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: "",
				ErrorCode:    errorCodes.MESSAGE_LIMIT_EXCEEDED,
			},
			textPayload,
		)
	}

	c.Locals("textPayload", textPayload)
	return c.Next()
}

func Text(c *fiber.Ctx) error {
	textPayload := c.Locals("textPayload").([]text.TextPayloadClient)

	var msgSlice []models.Message
	clientIdentifier := c.Locals("clientIdentifier").(string)
	waAccountIdentifier := c.Locals("waAccountIdentifier").(string)

	for _, msg := range textPayload {
		textMessage, err := makeTextMessage.MakeTextMessage(
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
				textPayload,
			)
		}
		msgSlice = append(msgSlice, *textMessage)
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
			textPayload,
		)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"ids": objIds,
	}, textPayload)
}
