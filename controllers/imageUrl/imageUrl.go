package imageUrl

import (
	"multiBot/config"
	"multiBot/errorCodes"
	"multiBot/models"
	makeImageUrlMessage "multiBot/sender/imageUrl"
	"multiBot/types/imageUrl"
	"multiBot/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

func ValiadateImageUrlPayload(c *fiber.Ctx) error {
	var imageUrlPayload []imageUrl.ImageUrlPayloadClient
	err := c.BodyParser(&imageUrlPayload)
	if err != nil {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: err.Error(),
				ErrorCode:    errorCodes.REQUEST_BODY_PARSING_ERROR,
			},
			imageUrlPayload,
		)
	}

	for _, msg := range imageUrlPayload {
		v := validate.Struct(msg)
		if !v.Validate() {
			return utils.ErrorResponse(
				c,
				utils.AppError{
					DebugMessage: v.Errors.One(),
					ErrorCode:    errorCodes.VALIDATION_ERROR,
				},
				imageUrlPayload,
			)
		}
	}

	cfg := config.GetConfig()
	messageLimit := cfg.API_MESSAGE_LIMIT
	if imageUrlPayload != nil && len(imageUrlPayload) > messageLimit {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: "",
				ErrorCode:    errorCodes.MESSAGE_LIMIT_EXCEEDED,
			},
			imageUrlPayload,
		)
	}

	c.Locals("imageUrlPayload", imageUrlPayload)
	return c.Next()
}

func ImageUrl(c *fiber.Ctx) error {
	imageUrlPayload := c.Locals("imageUrlPayload").([]imageUrl.ImageUrlPayloadClient)

	var msgSlice []models.Message
	clientIdentifier := c.Locals("clientIdentifier").(string)
	waAccountIdentifier := c.Locals("waAccountIdentifier").(string)

	for _, msg := range imageUrlPayload {
		imageUrlMessage, err := makeImageUrlMessage.MakeImageUrlMessage(
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
				imageUrlPayload,
			)
		}
		msgSlice = append(msgSlice, *imageUrlMessage)
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
			imageUrlPayload,
		)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"ids": objIds,
	}, imageUrlPayload)
}
