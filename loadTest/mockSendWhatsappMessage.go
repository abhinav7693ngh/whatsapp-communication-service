package loadTest

import (
	"multiBot/constants"
	"multiBot/logger"
	"multiBot/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func SendWhatsappMessageMock(message models.Message) {
	time.Sleep(800 * time.Millisecond)
	filter := bson.M{
		"_id": message.Id,
	}
	currentTime := time.Now().UnixMilli()
	update := bson.M{
		"$set": bson.M{
			"msgId":     "load_test",
			"status":    constants.MSG_STATUS_WA_PUSH_SUCCESS,
			"updatedAt": currentTime,
		},
		"$push": bson.M{
			"statusHistory": bson.M{
				"status":    constants.MSG_STATUS_WA_PUSH_SUCCESS,
				"updatedAt": currentTime,
			},
		},
		"$inc": bson.M{
			"retryCount": 1,
		},
	}

	updatedNo, err := message.UpdateOneMessage(filter, update)
	if err != nil {
		logger.LogDebug(
			nil,
			"SendWhatsappMessage_LOAD_TEST_MOCK: Error in updating message in status OK case with id: "+message.Id.Hex()+", error: "+err.Error(),
			nil,
		)
	} else if updatedNo != nil {
		if *updatedNo > 0 {
			logger.LogDebug(
				nil,
				"SendWhatsappMessage_LOAD_TEST_MOCK: Updated message in status OK case with id: "+message.Id.Hex(),
				nil,
			)
		}
	}
}
