package audioUrl

import (
	"multiBot/config"
	"multiBot/errorCodes"
	"multiBot/models"
	makeAudioUrlMessage "multiBot/sender/audioUrl"
	"multiBot/types/audioUrl"
	"multiBot/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

func ValiadateAudioUrlPayload(c *fiber.Ctx) error {
	var audioUrlPayload []audioUrl.AudioUrlPayloadClient
	err := c.BodyParser(&audioUrlPayload)
	if err != nil {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: err.Error(),
				ErrorCode:    errorCodes.REQUEST_BODY_PARSING_ERROR,
			},
			audioUrlPayload,
		)
	}

	for _, msg := range audioUrlPayload {
		v := validate.Struct(msg)
		if !v.Validate() {
			return utils.ErrorResponse(
				c,
				utils.AppError{
					DebugMessage: v.Errors.One(),
					ErrorCode:    errorCodes.VALIDATION_ERROR,
				},
				audioUrlPayload,
			)
		}
	}

	cfg := config.GetConfig()
	messageLimit := cfg.API_MESSAGE_LIMIT
	if audioUrlPayload != nil && len(audioUrlPayload) > messageLimit {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: "",
				ErrorCode:    errorCodes.MESSAGE_LIMIT_EXCEEDED,
			},
			audioUrlPayload,
		)
	}

	c.Locals("audioUrlPayload", audioUrlPayload)
	return c.Next()
}

func AudioUrl(c *fiber.Ctx) error {
	audioUrlPayload := c.Locals("audioUrlPayload").([]audioUrl.AudioUrlPayloadClient)

	var msgSlice []models.Message
	clientIdentifier := c.Locals("clientIdentifier").(string)
	waAccountIdentifier := c.Locals("waAccountIdentifier").(string)

	for _, msg := range audioUrlPayload {
		audioUrlMessage, err := makeAudioUrlMessage.MakeAudioUrlMessage(
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
				audioUrlPayload,
			)
		}
		msgSlice = append(msgSlice, *audioUrlMessage)
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
			audioUrlPayload,
		)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"ids": objIds,
	}, audioUrlPayload)
}
