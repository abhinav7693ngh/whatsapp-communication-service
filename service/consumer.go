package service

import (
	"context"
	"fmt"
	"multiBot/config"
	"multiBot/constants"
	"multiBot/logger"
	"multiBot/models"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConsumerWhatsapp(exitCtx context.Context, shutDownWg *sync.WaitGroup) {
	defer shutDownWg.Done()
	fmt.Println("Cosumer whatsapp process running !")

	cfg := config.GetConfig()

	messageLimit := cfg.WHATSAPP_CONSUMER.GET_MESSAGE_LIMIT
	messageFrequency := cfg.WHATSAPP_CONSUMER.GET_MESSAGE_FREQUENCY_SECONDS
	maxRetryAvailable := cfg.WHATSAPP_CONSUMER.MAX_WHATSAPP_PUSH_RETRY_AVAILABLE

	getMsgsForProcessing := []string{
		string(constants.MSG_STATUS_CREATED),
		string(constants.MSG_STATUS_WA_PUSH_FAILED),
		// TODO: Need to add timeout case as well, Do Separate handling of that
	}

	for {
		select {
		case <-exitCtx.Done():
			fmt.Println("Shutting down consumer whatsapp go routine !")
			return
		case <-time.After(time.Second * messageFrequency):
			{
				func() {
					defer func() {
						if r := recover(); r != nil {
							errStr := fmt.Sprint(r)
							errMsg := "Recovery: Panic occurred in ConsumerWhatsapp GoRoutine" + ", error: " + errStr
							logger.LogPanic(nil, errMsg, nil)
						}
					}()

					// // =========== FOR LOAD TEST ================= //

					// start := time.Now()

					// noOfMsg := 0

					// // =========================================== //

					currentTime := time.Now().UnixMilli()
					findFilter := bson.M{
						"status": bson.M{
							"$in": getMsgsForProcessing,
						},
						"retryCount": bson.M{
							"$lt": maxRetryAvailable,
						},
					}
					after := options.After
					findUpdateOptions := options.FindOneAndUpdateOptions{
						ReturnDocument: &after,
						Sort:           bson.M{"createdAt": 1},
					}
					updateOptions := bson.M{
						"$set": bson.M{
							"status":    constants.MSG_STATUS_PROCESSING,
							"updatedAt": currentTime,
						},
						"$push": bson.M{
							"statusHistory": bson.M{
								"status":    constants.MSG_STATUS_PROCESSING,
								"updatedAt": currentTime,
							},
						},
					}
					var messageModel models.Message
					wg := sync.WaitGroup{}
					for i := 0; i < int(messageLimit); i++ {
						message, err := messageModel.FindOneAndUpdate(findFilter, updateOptions, &findUpdateOptions)
						if err != nil {
							continue
						} else if message != nil {
							// // ============ FOR LOAD TEST ============= //
							// noOfMsg += 1
							// // ======================================== //
							logger.LogInfo(
								nil,
								"ConsumerWhatsapp GoRoutine: Updated msg status to processing in db of msg with id: "+message.Id.Hex(),
								nil,
							)
							wg.Add(1)
							go func(msg models.Message) {
								defer wg.Done()
								provider := NotificationProviderWhatsapp{}
								provider.SendWhatsappMessage(msg)

								// // ============ FOR LOAD TEST ============= //
								// loadTest.SendWhatsappMessageMock(msg)
								// // ======================================== //
							}(*message)
						}
					}
					wg.Wait()

					// // =========== FOR LOAD TEST ================= //

					// if noOfMsg > 0 {
					// 	str := "WA_CONSUMER_LOAD_TEST: Time taken to process " + strconv.FormatInt(messageLimit, 10) + " messages with frequency " + strconv.FormatInt(int64(messageFrequency), 10) + " seconds, updating " + strconv.FormatInt(int64(noOfMsg), 10) + " messages"
					// 	elapsed := time.Since(start)
					// 	str = str + " is " + elapsed.String()
					// 	logger.LogDebug(nil, str, nil)
					// }

					// // ========================================== //

				}()
			}
		}
	}
}
