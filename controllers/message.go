package controllers

import (
	"encoding/json"
	"multiBot/config"
	"multiBot/constants"
	"multiBot/errorCodes"
	"multiBot/models"
	makeAudioUrlMesssage "multiBot/sender/audioUrl"
	makeDocumentUrlMessage "multiBot/sender/documentUrl"
	makeImageUrlMessage "multiBot/sender/imageUrl"
	makeListMessage "multiBot/sender/list"
	makeLocationMessage "multiBot/sender/location"
	makeReplyButtonMessage "multiBot/sender/replyButton"
	makeStickerUrlMessage "multiBot/sender/stickerUrl"
	makeTemplateMessage "multiBot/sender/template"
	makeTextMessage "multiBot/sender/text"
	makeVideoUrlMessage "multiBot/sender/videoUrl"
	"multiBot/types"
	"multiBot/types/audioUrl"
	"multiBot/types/documentUrl"
	"multiBot/types/imageUrl"
	"multiBot/types/list"
	"multiBot/types/location"
	"multiBot/types/replyButton"
	"multiBot/types/stickerUrl"
	"multiBot/types/template"
	"multiBot/types/text"
	"multiBot/types/videoUrl"
	"multiBot/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

func ValiadateMessage(c *fiber.Ctx) error {
	var messagePayload []types.MessagePayloadClient
	err := json.Unmarshal(c.Body(), &messagePayload)
	if err != nil {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: err.Error(),
				ErrorCode:    errorCodes.UNMARSHALING_ERROR,
			},
			messagePayload,
		)
	}

	for _, msg := range messagePayload {
		v := validate.Struct(msg)
		if !v.Validate() {
			return utils.ErrorResponse(
				c,
				utils.AppError{
					DebugMessage: v.Errors.One(),
					ErrorCode:    errorCodes.VALIDATION_ERROR,
				},
				messagePayload,
			)
		}
		switch msg.Type {
		case string(constants.MSG_TYPE_TEXT):
			{
				msgBody, ok := msg.Body.(text.TextDataStruct)
				if !ok {
					return utils.ErrorResponse(
						c,
						utils.AppError{
							DebugMessage: "",
							ErrorCode:    errorCodes.INVALID_BODY,
						},
						messagePayload,
					)
				}
				v := validate.Struct(msgBody)
				if !v.Validate() {
					return utils.ErrorResponse(
						c,
						utils.AppError{
							DebugMessage: v.Errors.One(),
							ErrorCode:    errorCodes.VALIDATION_ERROR,
						},
						messagePayload,
					)
				}
			}
		case string(constants.MSG_TYPE_TEMPLATE):
			{
				msgBody, ok := msg.Body.(template.TemplateDataStruct)
				if !ok {
					return utils.ErrorResponse(
						c,
						utils.AppError{
							DebugMessage: "",
							ErrorCode:    errorCodes.INVALID_BODY,
						},
						messagePayload,
					)
				}
				v := validate.Struct(msgBody)
				if !v.Validate() {
					return utils.ErrorResponse(
						c,
						utils.AppError{
							DebugMessage: v.Errors.One(),
							ErrorCode:    errorCodes.VALIDATION_ERROR,
						},
						messagePayload,
					)
				}
			}
		case string(constants.MSG_TYPE_IMAGE_URL):
			{
				msgBody, ok := msg.Body.(imageUrl.ImageUrlDataStruct)
				if !ok {
					return utils.ErrorResponse(
						c,
						utils.AppError{
							DebugMessage: "",
							ErrorCode:    errorCodes.INVALID_BODY,
						},
						messagePayload,
					)
				}
				v := validate.Struct(msgBody)
				if !v.Validate() {
					return utils.ErrorResponse(
						c,
						utils.AppError{
							DebugMessage: v.Errors.One(),
							ErrorCode:    errorCodes.VALIDATION_ERROR,
						},
						messagePayload,
					)
				}
			}
		case string(constants.MSG_TYPE_AUDIO_URL):
			{
				msgBody, ok := msg.Body.(audioUrl.AudioUrlDataStruct)
				if !ok {
					return utils.ErrorResponse(
						c,
						utils.AppError{
							DebugMessage: "",
							ErrorCode:    errorCodes.INVALID_BODY,
						},
						messagePayload,
					)
				}
				v := validate.Struct(msgBody)
				if !v.Validate() {
					return utils.ErrorResponse(
						c,
						utils.AppError{
							DebugMessage: v.Errors.One(),
							ErrorCode:    errorCodes.VALIDATION_ERROR,
						},
						messagePayload,
					)
				}
			}
		case string(constants.MSG_TYPE_DOCUMENT_URL):
			{
				msgBody, ok := msg.Body.(documentUrl.DocumentUrlDataStruct)
				if !ok {
					return utils.ErrorResponse(
						c,
						utils.AppError{
							DebugMessage: "",
							ErrorCode:    errorCodes.INVALID_BODY,
						},
						messagePayload,
					)
				}
				v := validate.Struct(msgBody)
				if !v.Validate() {
					return utils.ErrorResponse(
						c,
						utils.AppError{
							DebugMessage: v.Errors.One(),
							ErrorCode:    errorCodes.VALIDATION_ERROR,
						},
						messagePayload,
					)
				}
			}
		case string(constants.MSG_TYPE_VIDEO_URL):
			{
				msgBody, ok := msg.Body.(videoUrl.VideoUrlDataStruct)
				if !ok {
					return utils.ErrorResponse(
						c,
						utils.AppError{
							DebugMessage: "",
							ErrorCode:    errorCodes.INVALID_BODY,
						},
						messagePayload,
					)
				}
				v := validate.Struct(msgBody)
				if !v.Validate() {
					return utils.ErrorResponse(
						c,
						utils.AppError{
							DebugMessage: v.Errors.One(),
							ErrorCode:    errorCodes.VALIDATION_ERROR,
						},
						messagePayload,
					)
				}
			}
		case string(constants.MSG_TYPE_STICKER_URL):
			{
				msgBody, ok := msg.Body.(stickerUrl.StickerUrlDataStruct)
				if !ok {
					return utils.ErrorResponse(
						c,
						utils.AppError{
							DebugMessage: "",
							ErrorCode:    errorCodes.INVALID_BODY,
						},
						messagePayload,
					)
				}
				v := validate.Struct(msgBody)
				if !v.Validate() {
					return utils.ErrorResponse(
						c,
						utils.AppError{
							DebugMessage: v.Errors.One(),
							ErrorCode:    errorCodes.VALIDATION_ERROR,
						},
						messagePayload,
					)
				}
			}
		case string(constants.MSG_TYPE_LOCATION):
			{
				msgBody, ok := msg.Body.(location.LocationDataStruct)
				if !ok {
					return utils.ErrorResponse(
						c,
						utils.AppError{
							DebugMessage: "",
							ErrorCode:    errorCodes.INVALID_BODY,
						},
						messagePayload,
					)
				}
				v := validate.Struct(msgBody)
				if !v.Validate() {
					return utils.ErrorResponse(
						c,
						utils.AppError{
							DebugMessage: v.Errors.One(),
							ErrorCode:    errorCodes.VALIDATION_ERROR,
						},
						messagePayload,
					)
				}
			}
		case string(constants.MSG_TYPE_LIST_INTERNAL):
			{
				msgBody, ok := msg.Body.(list.ListDataStruct)
				if !ok {
					return utils.ErrorResponse(
						c,
						utils.AppError{
							DebugMessage: "",
							ErrorCode:    errorCodes.INVALID_BODY,
						},
						messagePayload,
					)
				}
				v := validate.Struct(msgBody)
				if !v.Validate() {
					return utils.ErrorResponse(
						c,
						utils.AppError{
							DebugMessage: v.Errors.One(),
							ErrorCode:    errorCodes.VALIDATION_ERROR,
						},
						messagePayload,
					)
				}
			}
		case string(constants.MSG_TYPE_REPLY_BUTTON_INTERNAL):
			{
				msgBody, ok := msg.Body.(replyButton.ReplyButtonDataStruct)
				if !ok {
					return utils.ErrorResponse(
						c,
						utils.AppError{
							DebugMessage: "",
							ErrorCode:    errorCodes.INVALID_BODY,
						},
						messagePayload,
					)
				}
				v := validate.Struct(msgBody)
				if !v.Validate() {
					return utils.ErrorResponse(
						c,
						utils.AppError{
							DebugMessage: v.Errors.One(),
							ErrorCode:    errorCodes.VALIDATION_ERROR,
						},
						messagePayload,
					)
				}
			}
		default:
			{
				return utils.ErrorResponse(
					c,
					utils.AppError{
						DebugMessage: v.Errors.One(),
						ErrorCode:    errorCodes.INVALID_TYPE,
					},
					messagePayload,
				)
			}
		}
	}

	cfg := config.GetConfig()
	messageLimit := cfg.API_MESSAGE_LIMIT
	if messagePayload != nil && len(messagePayload) > messageLimit {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: "",
				ErrorCode:    errorCodes.MESSAGE_LIMIT_EXCEEDED,
			},
			messagePayload,
		)
	}

	c.Locals("messagePayload", messagePayload)
	return c.Next()
}

func Message(c *fiber.Ctx) error {
	messagePayload := c.Locals("messagePayload").([]types.MessagePayloadClient)

	var msgSlice []models.Message
	clientIdentifier := c.Locals("clientIdentifier").(string)
	waAccountIdentifier := c.Locals("waAccountIdentifier").(string)

	for _, msg := range messagePayload {
		switch msg.Type {
		case string(constants.MSG_TYPE_TEXT):
			{
				message, err := makeTextMessage.MakeTextMessage(
					msg.Body.(text.TextDataStruct),
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
						messagePayload,
					)
				}
				msgSlice = append(msgSlice, *message)
			}
		case string(constants.MSG_TYPE_TEMPLATE):
			{
				message, err := makeTemplateMessage.MakeTemplateMessage(
					msg.Body.(template.TemplateDataStruct),
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
						messagePayload,
					)
				}
				msgSlice = append(msgSlice, *message)
			}
		case string(constants.MSG_TYPE_IMAGE_URL):
			{
				message, err := makeImageUrlMessage.MakeImageUrlMessage(
					msg.Body.(imageUrl.ImageUrlDataStruct),
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
						messagePayload,
					)
				}
				msgSlice = append(msgSlice, *message)
			}
		case string(constants.MSG_TYPE_AUDIO_URL):
			{
				message, err := makeAudioUrlMesssage.MakeAudioUrlMessage(
					msg.Body.(audioUrl.AudioUrlDataStruct),
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
						messagePayload,
					)
				}
				msgSlice = append(msgSlice, *message)
			}
		case string(constants.MSG_TYPE_DOCUMENT_URL):
			{
				message, err := makeDocumentUrlMessage.MakeDocumentUrlMessage(
					msg.Body.(documentUrl.DocumentUrlDataStruct),
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
						messagePayload,
					)
				}
				msgSlice = append(msgSlice, *message)
			}
		case string(constants.MSG_TYPE_VIDEO_URL):
			{
				message, err := makeVideoUrlMessage.MakeVideoUrlMessage(
					msg.Body.(videoUrl.VideoUrlDataStruct),
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
						messagePayload,
					)
				}
				msgSlice = append(msgSlice, *message)
			}
		case string(constants.MSG_TYPE_STICKER_URL):
			{
				message, err := makeStickerUrlMessage.MakeStickerUrlMessage(
					msg.Body.(stickerUrl.StickerUrlDataStruct),
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
						messagePayload,
					)
				}
				msgSlice = append(msgSlice, *message)
			}
		case string(constants.MSG_TYPE_LOCATION):
			{
				message, err := makeLocationMessage.MakeLocationMessage(
					msg.Body.(location.LocationDataStruct),
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
						messagePayload,
					)
				}
				msgSlice = append(msgSlice, *message)
			}
		case string(constants.MSG_TYPE_LIST_INTERNAL):
			{
				message, err := makeListMessage.MakeListMessage(
					msg.Body.(list.ListDataStruct),
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
						messagePayload,
					)
				}
				msgSlice = append(msgSlice, *message)
			}
		case string(constants.MSG_TYPE_REPLY_BUTTON_INTERNAL):
			{
				message, err := makeReplyButtonMessage.MakeReplyButtonMessage(
					msg.Body.(replyButton.ReplyButtonDataStruct),
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
						messagePayload,
					)
				}
				msgSlice = append(msgSlice, *message)
			}
		}
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
			messagePayload,
		)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"ids": objIds,
	}, messagePayload)
}
