package list

import (
	"multiBot/config"
	"multiBot/errorCodes"
	"multiBot/models"
	makeListMessage "multiBot/sender/list"
	"multiBot/types/list"
	"multiBot/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

func ValiadateListPayload(c *fiber.Ctx) error {
	var listPayload []list.ListPayloadClient
	err := c.BodyParser(&listPayload)
	if err != nil {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: err.Error(),
				ErrorCode:    errorCodes.REQUEST_BODY_PARSING_ERROR,
			},
			listPayload,
		)
	}

	for _, msg := range listPayload {
		v := validate.Struct(msg)
		if !v.Validate() {
			return utils.ErrorResponse(
				c,
				utils.AppError{
					DebugMessage: v.Errors.One(),
					ErrorCode:    errorCodes.VALIDATION_ERROR,
				},
				listPayload,
			)
		}
	}

	cfg := config.GetConfig()
	messageLimit := cfg.API_MESSAGE_LIMIT
	if listPayload != nil && len(listPayload) > messageLimit {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: "",
				ErrorCode:    errorCodes.MESSAGE_LIMIT_EXCEEDED,
			},
			listPayload,
		)
	}

	c.Locals("listPayload", listPayload)
	return c.Next()
}

func List(c *fiber.Ctx) error {
	listPayload := c.Locals("listPayload").([]list.ListPayloadClient)

	var msgSlice []models.Message
	clientIdentifier := c.Locals("clientIdentifier").(string)
	waAccountIdentifier := c.Locals("waAccountIdentifier").(string)

	for _, msg := range listPayload {
		listMessage, err := makeListMessage.MakeListMessage(
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
				listPayload,
			)
		}
		msgSlice = append(msgSlice, *listMessage)
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
			listPayload,
		)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"ids": objIds,
	}, listPayload)
}
