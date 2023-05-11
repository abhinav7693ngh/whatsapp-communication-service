package controllers

import (
	"multiBot/config"
	"multiBot/errorCodes"
	"multiBot/models"
	"multiBot/types"
	"multiBot/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValiadateStatusPayload(c *fiber.Ctx) error {
	var statusPayload types.StatusPayload
	err := c.BodyParser(&statusPayload)
	if err != nil {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: err.Error(),
				ErrorCode:    errorCodes.REQUEST_BODY_PARSING_ERROR,
			},
			statusPayload,
		)
	}

	v := validate.Struct(statusPayload)
	if !v.Validate() {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: v.Errors.One(),
				ErrorCode:    errorCodes.VALIDATION_ERROR,
			},
			statusPayload,
		)
	}

	cfg := config.GetConfig()
	messageLimit := cfg.API_MESSAGE_LIMIT
	if statusPayload.Ids != nil && len(statusPayload.Ids) > messageLimit {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: "",
				ErrorCode:    errorCodes.MESSAGE_LIMIT_EXCEEDED,
			},
			statusPayload,
		)
	}

	c.Locals("statusPayload", statusPayload)
	return c.Next()
}

func Status(c *fiber.Ctx) error {
	statusPayload := c.Locals("statusPayload").(types.StatusPayload)

	if len(statusPayload.Ids) == 0 {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: "",
				ErrorCode:    errorCodes.NO_MESSAGE_ID_PROVIDED,
			},
			statusPayload,
		)
	}

	var objectIds []primitive.ObjectID
	for _, id := range statusPayload.Ids {
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return utils.ErrorResponse(
				c,
				utils.AppError{
					DebugMessage: err.Error(),
					ErrorCode:    errorCodes.DB_OBJECT_ID_CONVERSION_ERROR,
				},
				statusPayload,
			)
		}
		objectIds = append(objectIds, objectId)
	}

	filter := bson.M{"_id": bson.M{"$in": objectIds}}
	var messageModel models.Message
	msgs, err := messageModel.FindManyMessages(filter, nil)
	if err != nil {
		if err.Error() == errorCodes.NO_RECORD_FOUND {
			return utils.SuccessResponse(c, types.StatusResponse{
				Count:    0,
				Messages: nil,
			}, statusPayload)
		}
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: err.Error(),
				ErrorCode:    errorCodes.DB_OPERATION_ERROR,
			},
			statusPayload,
		)
	}

	var statusResp []types.StatusMessage
	for _, msg := range *msgs {
		waAccountInfo := config.GetWhatsappAccountInfo(msg.WhatsappAccountIdentifier)
		clientInfo := config.GetClientInfo(msg.ClientIdentifier)
		statusResp = append(statusResp, types.StatusMessage{
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

	return utils.SuccessResponse(c, types.StatusResponse{
		Count:    len(statusResp),
		Messages: statusResp,
	}, statusPayload)
}
