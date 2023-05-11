package service

import (
	"context"
	"fmt"
	"multiBot/config"
	"multiBot/constants"
	"multiBot/logger"
	"multiBot/models"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdateStaleMessages(exitCtx context.Context, shutDownWg *sync.WaitGroup) {
	defer shutDownWg.Done()
	fmt.Println("Update stale messages process running !")

	cfg := config.GetConfig()

	staleMessageFrequency := cfg.WHATSAPP_CONSUMER.UPDATE_STALE_MESSAGES_FREQUENCY_SECONDS

	maxRetryAvailable := cfg.WHATSAPP_CONSUMER.MAX_WHATSAPP_PUSH_RETRY_AVAILABLE

	getStaleMessages := []string{string(constants.MSG_STATUS_PROCESSING)}

	for {
		select {
		case <-exitCtx.Done():
			fmt.Println("Shutting down update stale messages go routine !")
			return
		case <-time.After(time.Second * staleMessageFrequency):
			{
				func() {
					defer func() {
						if r := recover(); r != nil {
							errStr := fmt.Sprint(r)
							errMsg := "Recovery: Panic occurred in UpdateStaleMessages GoRoutine" + ", error: " + errStr
							logger.LogPanic(nil, errMsg, nil)
						}
					}()

					// older than 2 minutes messages with status processing filter with retryCount less than max available
					minutesAgo := time.Now().Add(-2 * time.Minute).UnixMilli()
					filter := bson.M{
						"status": bson.M{
							"$in": getStaleMessages,
						},
						"updatedAt": bson.M{
							"$lt": minutesAgo,
						},
						"retryCount": bson.M{
							"$lt": maxRetryAvailable,
						},
					}

					// Update to status failed and increased the retryCount
					currentTime := time.Now().UnixMilli()
					update := bson.M{
						"$set": bson.M{
							"status":    constants.MSG_STATUS_WA_PUSH_FAILED,
							"updatedAt": currentTime,
						},
						"$push": bson.M{
							"statusHistory": bson.M{
								"status":    constants.MSG_STATUS_WA_PUSH_FAILED,
								"updatedAt": currentTime,
							},
						},
						"$inc": bson.M{
							"retryCount": 1,
						},
					}

					var message models.Message
					updatedNo, err := message.UpdateManyMessages(filter, update)
					if err != nil {
						logger.LogError(
							nil,
							"UpdateStaleMessages GoRoutine: Error in updating stale messages status, error: "+err.Error(),
							nil,
						)
					} else if updatedNo != nil {
						if *updatedNo > 0 {
							logger.LogInfo(
								nil,
								"UpdateStaleMessages GoRoutine: Updated stale messages status, number of updated messages: "+strconv.FormatInt(*updatedNo, 10),
								nil,
							)
						}
					}
				}()
			}
		}
	}
}
