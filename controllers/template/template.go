package template

import (
	"multiBot/config"
	"multiBot/errorCodes"
	"multiBot/models"
	makeTemplateMessage "multiBot/sender/template"
	"multiBot/types/template"
	"multiBot/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

func ValiadateTemplatePayload(c *fiber.Ctx) error {
	var templatePayload []template.TemplatePayloadClient
	err := c.BodyParser(&templatePayload)
	if err != nil {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: err.Error(),
				ErrorCode:    errorCodes.REQUEST_BODY_PARSING_ERROR,
			},
			templatePayload,
		)
	}

	for _, msg := range templatePayload {
		v := validate.Struct(msg)
		if !v.Validate() {
			return utils.ErrorResponse(
				c,
				utils.AppError{
					DebugMessage: v.Errors.One(),
					ErrorCode:    errorCodes.VALIDATION_ERROR,
				},
				templatePayload,
			)
		}
	}

	cfg := config.GetConfig()
	messageLimit := cfg.API_MESSAGE_LIMIT
	if templatePayload != nil && len(templatePayload) > messageLimit {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: "",
				ErrorCode:    errorCodes.MESSAGE_LIMIT_EXCEEDED,
			},
			templatePayload,
		)
	}

	c.Locals("templatePayload", templatePayload)
	return c.Next()
}

func Template(c *fiber.Ctx) error {
	templatePayload := c.Locals("templatePayload").([]template.TemplatePayloadClient)

	var msgSlice []models.Message
	clientIdentifier := c.Locals("clientIdentifier").(string)
	waAccountIdentifier := c.Locals("waAccountIdentifier").(string)

	for _, msg := range templatePayload {
		templateMessage, err := makeTemplateMessage.MakeTemplateMessage(
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
				templatePayload,
			)
		}
		msgSlice = append(msgSlice, *templateMessage)
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
			templatePayload,
		)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"ids": objIds,
	}, templatePayload)
}
