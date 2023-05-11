package queue

import (
	"context"
	"fmt"
	"multiBot/config"
	"multiBot/logger"
	"multiBot/service"
	"strings"
	"sync"

	kafkaGo "github.com/segmentio/kafka-go"
)

func Consumer(exitCtx context.Context, shutDownWg *sync.WaitGroup) {
	defer shutDownWg.Done()
	cfg := config.GetConfig()

	groupId := cfg.KAFKA.GROUP_ID
	brokersAddress := strings.Split(cfg.KAFKA.BROKERS_ADDR, ",")
	topic := cfg.KAFKA.TOPIC

	kafkaConsumer := kafkaGo.NewReader(kafkaGo.ReaderConfig{
		Brokers:  brokersAddress,
		GroupID:  groupId,
		Topic:    topic,
		MinBytes: 1,    // 1B
		MaxBytes: 10e6, // 10MB
	})

	fmt.Println("Consumer process running !")

	for {
		select {
		case <-exitCtx.Done():
			{
				if err := kafkaConsumer.Close(); err != nil {
					logger.LogError(nil, "Failed to close consumer: "+err.Error(), nil)
				}
				fmt.Println("Shutting down queue consumer !")
				return
			}
		default:
			{
				func() {
					defer func() {
						if r := recover(); r != nil {
							errStr := fmt.Sprint(r)
							errMsg := "Recovery: Panic occurred in queue Consumer GoRoutine" + ", error: " + errStr
							logger.LogPanic(nil, errMsg, nil)
						}
					}()
					msg, err := kafkaConsumer.ReadMessage(exitCtx)
					// TODO: Check if ReadMessage, retry if failed to consume message

					if err != nil {
						logger.LogError(
							nil,
							"Consumer GoRoutine: Failed to consume message: "+err.Error()+", msgId: "+string(msg.Key),
							nil,
						)
					} else {
						logger.LogInfo(
							nil,
							"Consumer GoRoutine: Message successfully consumed, msgId: "+string(msg.Key)+", status: "+string(msg.Value),
							nil,
						)
						service.UpdateWebhookMsgStatus(string(msg.Key), string(msg.Value))
					}
				}()
			}
		}
	}
}
