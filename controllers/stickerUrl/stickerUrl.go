package stickerUrl

import (
	"multiBot/config"
	"multiBot/errorCodes"
	"multiBot/models"
	makeStickerUrlMessage "multiBot/sender/stickerUrl"
	"multiBot/types/stickerUrl"
	"multiBot/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

func ValiadateStickerUrlPayload(c *fiber.Ctx) error {
	var stickerUrlPayload []stickerUrl.StickerUrlPayloadClient
	err := c.BodyParser(&stickerUrlPayload)
	if err != nil {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: err.Error(),
				ErrorCode:    errorCodes.REQUEST_BODY_PARSING_ERROR,
			},
			stickerUrlPayload,
		)
	}

	for _, msg := range stickerUrlPayload {
		v := validate.Struct(msg)
		if !v.Validate() {
			return utils.ErrorResponse(
				c,
				utils.AppError{
					DebugMessage: v.Errors.One(),
					ErrorCode:    errorCodes.VALIDATION_ERROR,
				},
				stickerUrlPayload,
			)
		}
	}

	cfg := config.GetConfig()
	messageLimit := cfg.API_MESSAGE_LIMIT
	if stickerUrlPayload != nil && len(stickerUrlPayload) > messageLimit {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: "",
				ErrorCode:    errorCodes.MESSAGE_LIMIT_EXCEEDED,
			},
			stickerUrlPayload,
		)
	}

	c.Locals("stickerUrlPayload", stickerUrlPayload)
	return c.Next()
}

func StickerUrl(c *fiber.Ctx) error {
	stickerUrlPayload := c.Locals("stickerUrlPayload").([]stickerUrl.StickerUrlPayloadClient)

	var msgSlice []models.Message
	clientIdentifier := c.Locals("clientIdentifier").(string)
	waAccountIdentifier := c.Locals("waAccountIdentifier").(string)

	for _, msg := range stickerUrlPayload {
		stickerUrlMessage, err := makeStickerUrlMessage.MakeStickerUrlMessage(
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
				stickerUrlPayload,
			)
		}
		msgSlice = append(msgSlice, *stickerUrlMessage)
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
			stickerUrlPayload,
		)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"ids": objIds,
	}, stickerUrlPayload)
}
