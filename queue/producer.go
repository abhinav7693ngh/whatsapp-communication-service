package queue

import (
	"context"
	"fmt"
	"multiBot/config"
	"multiBot/constants"
	"multiBot/logger"
	"multiBot/models"
	"strings"
	"sync"
	"time"

	kafkaGo "github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func updateMessageStatus(msgId string, status string, done bool) {
	var messageModel models.Message
	filter := bson.M{
		"msgId": msgId,
	}
	var update primitive.M
	if done {
		possible := models.PossibleToUpdateStatusWithFilter(constants.MSG_STATUS_WA_STATUS_RECEIVED_QUEUE_PUSH_SUCCESS, filter, msgId, "updateMessageStatus")
		if !possible {
			return
		}
		currentTime := time.Now().UnixMilli()
		update = bson.M{
			"$set": bson.M{
				"status":    constants.MSG_STATUS_WA_STATUS_RECEIVED_QUEUE_PUSH_SUCCESS,
				"updatedAt": currentTime,
			},
			"$push": bson.M{
				"statusHistory": bson.M{
					"status":    constants.MSG_STATUS_WA_STATUS_RECEIVED_QUEUE_PUSH_SUCCESS,
					"updatedAt": currentTime,
				},
			},
		}
	} else {
		possible := models.PossibleToUpdateStatusWithFilter(constants.MSG_STATUS_WA_STATUS_RECEIVED_QUEUE_PUSH_FAILED, filter, msgId, "updateMessageStatus")
		if !possible {
			return
		}
		currentTime := time.Now().UnixMilli()
		update = bson.M{
			"$set": bson.M{
				"status":    constants.MSG_STATUS_WA_STATUS_RECEIVED_QUEUE_PUSH_FAILED,
				"updatedAt": currentTime,
			},
			"$push": bson.M{
				"statusHistory": bson.M{
					"status":    constants.MSG_STATUS_WA_STATUS_RECEIVED_QUEUE_PUSH_FAILED,
					"updatedAt": currentTime,
				},
			},
		}
	}
	updatedNo, err := messageModel.UpdateOneMessage(filter, update)
	if err != nil {
		logger.LogError(
			nil,
			"Producer GoRoutine: Not able to update msg status in db in producer ack with msgId: "+msgId+" and msgStatus: "+status+", error: "+err.Error(),
			nil,
		)
	} else if updatedNo != nil && *updatedNo > 0 {
		logger.LogInfo(
			nil,
			"Producer GoRoutine: Updated msg status in db in producer ack with msgId: "+msgId+" and msgStatus: "+status,
			nil,
		)
	}
}

func Producer(exitCtx context.Context, shutDownWg *sync.WaitGroup) {
	defer shutDownWg.Done()
	cfg := config.GetConfig()

	brokersAddress := strings.Split(cfg.KAFKA.BROKERS_ADDR, ",")
	topic := cfg.KAFKA.TOPIC

	kafkaProducer := &kafkaGo.Writer{
		Addr:  kafkaGo.TCP(brokersAddress...),
		Topic: topic,
	}

	fmt.Println("Producer process running !")

	for {
		select {
		case <-exitCtx.Done():
			{
				if err := kafkaProducer.Close(); err != nil {
					logger.LogError(nil, "Failed to close producer: "+err.Error(), nil)
				}
				fmt.Println("Shutting down queue producer !")
				return
			}
		case msg := <-MessageChannel:
			{
				func() {
					defer func() {
						if r := recover(); r != nil {
							errStr := fmt.Sprint(r)
							errMsg := "Recovery: Panic occurred in queue Producer GoRoutine" + ", error: " + errStr
							logger.LogPanic(nil, errMsg, nil)
						}
					}()
					err := kafkaProducer.WriteMessages(context.Background(), msg)
					if err != nil {
						logger.LogError(
							nil,
							"Producer GoRoutine: Failed to produce messages: "+err.Error()+" ,msgId: "+string(msg.Key)+" ,status: "+string(msg.Value),
							nil,
						)
						updateMessageStatus(string(msg.Key), string(msg.Value), false)
					} else {
						logger.LogInfo(
							nil,
							"Producer GoRoutine: Message successfully produced, msgId: "+string(msg.Key)+", status: "+string(msg.Value),
							nil,
						)
						updateMessageStatus(string(msg.Key), string(msg.Value), true)
					}
				}()
			}
		}
	}
}
