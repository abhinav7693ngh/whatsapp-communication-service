package service

import (
	"multiBot/constants"
	"multiBot/logger"
	"multiBot/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdateWebhookMsgStatus(msgId string, status string) {
	var messageModel models.Message
	filter := bson.M{"msgId": msgId}
	var update primitive.M
	currentTime := time.Now().UnixMilli()
	if constants.CheckWhatsappStatus(status) {
		possible := models.PossibleToUpdateStatusWithFilter(constants.WebhookStatusMap[status], filter, msgId, "UpdateWebhookMsgStatus")
		if !possible {
			return
		}
		update = bson.M{
			"$set": bson.M{
				"status":    constants.WebhookStatusMap[status],
				"updatedAt": currentTime,
			},
			"$push": bson.M{
				"statusHistory": bson.M{
					"status":    constants.WebhookStatusMap[status],
					"updatedAt": currentTime,
				},
			},
		}

		updatedNo, err := messageModel.UpdateOneMessage(filter, update)
		if err != nil {
			logger.LogError(
				nil,
				"UpdateWebhookMsgStatus: Not able to update msg status in db in WA webhook ack with msgId: "+msgId+" and msgStatus: "+status+", error: "+err.Error(),
				nil,
			)
		} else if updatedNo != nil {
			if *updatedNo > 0 {
				logger.LogInfo(
					nil,
					"UpdateWebhookMsgStatus: Updated msg status in db in WA webhook ack with msgId: "+msgId+" and msgStatus: "+status,
					nil,
				)
			}
		} else {
			logger.LogInfo(
				nil,
				"UpdateWebhookMsgStatus: No msg updated in db in WA webhook ack with msgId: "+msgId+" and msgStatus: "+status,
				nil,
			)
		}
	} else {
		logger.LogError(
			nil,
			"UpdateWebhookMsgStatus: Got a different status message from whatsapp apart from sent, delivered, undelivered, failed or read in msgId: "+msgId+" and msgStatus: "+status,
			nil,
		)
	}
}
