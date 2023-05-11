package videoUrl

import (
	"multiBot/config"
	"multiBot/errorCodes"
	"multiBot/models"
	makeVideoUrlMessage "multiBot/sender/videoUrl"
	"multiBot/types/videoUrl"
	"multiBot/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

func ValiadateVideoUrlPayload(c *fiber.Ctx) error {
	var videoUrlPayload []videoUrl.VideoUrlPayloadClient
	err := c.BodyParser(&videoUrlPayload)
	if err != nil {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: err.Error(),
				ErrorCode:    errorCodes.REQUEST_BODY_PARSING_ERROR,
			},
			videoUrlPayload,
		)
	}

	for _, msg := range videoUrlPayload {
		v := validate.Struct(msg)
		if !v.Validate() {
			return utils.ErrorResponse(
				c,
				utils.AppError{
					DebugMessage: v.Errors.One(),
					ErrorCode:    errorCodes.VALIDATION_ERROR,
				},
				videoUrlPayload,
			)
		}
	}

	cfg := config.GetConfig()
	messageLimit := cfg.API_MESSAGE_LIMIT
	if videoUrlPayload != nil && len(videoUrlPayload) > messageLimit {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: "",
				ErrorCode:    errorCodes.MESSAGE_LIMIT_EXCEEDED,
			},
			videoUrlPayload,
		)
	}

	c.Locals("videoUrlPayload", videoUrlPayload)
	return c.Next()
}

func VideoUrl(c *fiber.Ctx) error {
	videoUrlPayload := c.Locals("videoUrlPayload").([]videoUrl.VideoUrlPayloadClient)

	var msgSlice []models.Message
	clientIdentifier := c.Locals("clientIdentifier").(string)
	waAccountIdentifier := c.Locals("waAccountIdentifier").(string)

	for _, msg := range videoUrlPayload {
		videoUrlMessage, err := makeVideoUrlMessage.MakeVideoUrlMessage(
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
				videoUrlPayload,
			)
		}
		msgSlice = append(msgSlice, *videoUrlMessage)
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
			videoUrlPayload,
		)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"ids": objIds,
	}, videoUrlPayload)
}
