package controllers

import (
	"multiBot/config"
	"multiBot/constants"
	"multiBot/errorCodes"
	"multiBot/models"
	"multiBot/types"
	"multiBot/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ValiadateGetMessagesUsingStatusPayload(c *fiber.Ctx) error {
	var getMessagesPayload types.GetMessagesUsingStatusPayload
	err := c.BodyParser(&getMessagesPayload)
	if err != nil {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: err.Error(),
				ErrorCode:    errorCodes.REQUEST_BODY_PARSING_ERROR,
			},
			getMessagesPayload,
		)
	}

	v := validate.Struct(getMessagesPayload)
	if !v.Validate() {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: v.Errors.One(),
				ErrorCode:    errorCodes.VALIDATION_ERROR,
			},
			getMessagesPayload,
		)
	}

	if !constants.CheckValidStatus(getMessagesPayload.Status) {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: v.Errors.One(),
				ErrorCode:    errorCodes.INVALID_MESSAGE_STATUS,
			},
			getMessagesPayload,
		)
	}

	c.Locals("getMessagesPayload", getMessagesPayload)
	return c.Next()
}

func GetMessagesUsingStatus(c *fiber.Ctx) error {
	getMessagesPayload := c.Locals("getMessagesPayload").(types.GetMessagesUsingStatusPayload)

	// filter
	filter := bson.M{"status": getMessagesPayload.Status}

	// messageModel instance
	var messageModel models.Message
	var pageInt, limitInt int
	var totalDocCount *int64
	var findManyMessagesErr error

	// page and limit query parameters
	page := c.Query("page")
	limit := c.Query("limit")

	var msgs *[]models.Message

	// Set up options for MongoDB query
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"createdAt": -1})

	totalDocCount, err := messageModel.CountDocuments(filter, nil)
	if err != nil {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: err.Error(),
				ErrorCode:    errorCodes.DB_OPERATION_ERROR,
			},
			getMessagesPayload,
		)
	}

	if page != "" && limit != "" {
		// Convert query params to ints
		pageInt, err = strconv.Atoi(page)
		if err != nil {
			return utils.ErrorResponse(
				c,
				utils.AppError{
					DebugMessage: err.Error(),
					ErrorCode:    errorCodes.INTEGER_CONVERSION_ERROR,
				},
				getMessagesPayload,
			)
		}

		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			return utils.ErrorResponse(
				c,
				utils.AppError{
					DebugMessage: err.Error(),
					ErrorCode:    errorCodes.INTEGER_CONVERSION_ERROR,
				},
				getMessagesPayload,
			)
		}

		// Calculate offset
		offset := (pageInt - 1) * limitInt

		findOptions.SetSkip(int64(offset))
		findOptions.SetLimit(int64(limitInt))

		msgs, findManyMessagesErr = messageModel.FindManyMessages(filter, findOptions)
	} else {
		pageInt = 1
		limitInt = 10

		// Calculate offset
		offset := (pageInt - 1) * limitInt

		findOptions.SetSkip(int64(offset))
		findOptions.SetLimit(int64(limitInt))

		msgs, findManyMessagesErr = messageModel.FindManyMessages(filter, findOptions)
	}

	if findManyMessagesErr != nil {
		if err.Error() == errorCodes.NO_RECORD_FOUND {
			return utils.SuccessResponse(c, types.GetMessagesUsingStatusResponse{
				Page:     pageInt,
				Limit:    limitInt,
				Count:    0,
				Total:    *totalDocCount,
				Messages: nil,
			}, getMessagesPayload)
		}
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: err.Error(),
				ErrorCode:    errorCodes.DB_OPERATION_ERROR,
			},
			getMessagesPayload,
		)
	}
	var getMessageStatusResp []types.MessageUsingStatus
	for _, msg := range *msgs {
		waAccountInfo := config.GetWhatsappAccountInfo(msg.WhatsappAccountIdentifier)
		clientInfo := config.GetClientInfo(msg.ClientIdentifier)
		getMessageStatusResp = append(getMessageStatusResp, types.MessageUsingStatus{
			Id:                models.GetStringObjectIdFromBson(msg.Id),
			Body:              msg.Body,
			Status:            string(msg.Status),
			MsgId:             msg.MsgId,
			StatusHistory:     msg.StatusHistory,
			RetryCount:        msg.RetryCount,
			MaxRetryAvailable: msg.MaxRetryAvailable,
			Client:            clientInfo.NAME,
			Organization:      clientInfo.ORGANIZATION,
			WhatsappAccount:   waAccountInfo.NAME,
			NetworkType:       msg.MsgNetworkType,
			Purpose:           msg.MsgPurpose,
			CreatedAt:         msg.CreatedAt,
			UpdatedAt:         msg.UpdatedAt,
		})
	}
	return utils.SuccessResponse(c, types.GetMessagesUsingStatusResponse{
		Page:     pageInt,
		Limit:    limitInt,
		Count:    len(getMessageStatusResp),
		Total:    *totalDocCount,
		Messages: getMessageStatusResp,
	}, getMessagesPayload)
}
