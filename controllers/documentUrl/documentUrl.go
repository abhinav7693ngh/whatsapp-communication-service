package documentUrl

import (
	"multiBot/config"
	"multiBot/errorCodes"
	"multiBot/models"
	makeDocumentUrlMessage "multiBot/sender/documentUrl"
	"multiBot/types/documentUrl"
	"multiBot/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

func ValiadateDocumentUrlPayload(c *fiber.Ctx) error {
	var documentUrlPayload []documentUrl.DocumentUrlPayloadClient
	err := c.BodyParser(&documentUrlPayload)
	if err != nil {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: err.Error(),
				ErrorCode:    errorCodes.REQUEST_BODY_PARSING_ERROR,
			},
			documentUrlPayload,
		)
	}

	for _, msg := range documentUrlPayload {
		v := validate.Struct(msg)
		if !v.Validate() {
			return utils.ErrorResponse(
				c,
				utils.AppError{
					DebugMessage: v.Errors.One(),
					ErrorCode:    errorCodes.VALIDATION_ERROR,
				},
				documentUrlPayload,
			)
		}
	}

	cfg := config.GetConfig()
	messageLimit := cfg.API_MESSAGE_LIMIT
	if documentUrlPayload != nil && len(documentUrlPayload) > messageLimit {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: "",
				ErrorCode:    errorCodes.MESSAGE_LIMIT_EXCEEDED,
			},
			documentUrlPayload,
		)
	}

	c.Locals("documentUrlPayload", documentUrlPayload)
	return c.Next()
}

func DocumentUrl(c *fiber.Ctx) error {
	documentUrlPayload := c.Locals("documentUrlPayload").([]documentUrl.DocumentUrlPayloadClient)

	var msgSlice []models.Message
	clientIdentifier := c.Locals("clientIdentifier").(string)
	waAccountIdentifier := c.Locals("waAccountIdentifier").(string)

	for _, msg := range documentUrlPayload {
		documentUrlMessage, err := makeDocumentUrlMessage.MakeDocumentUrlMessage(
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
				documentUrlPayload,
			)
		}
		msgSlice = append(msgSlice, *documentUrlMessage)
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
			documentUrlPayload,
		)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"ids": objIds,
	}, documentUrlPayload)
}
