package webhook

import (
	"multiBot/config"
	"multiBot/constants"
	"multiBot/logger"
	"multiBot/queue"
	textReplier "multiBot/replier/text"
	textSender "multiBot/sender/text"
	"multiBot/service"
	"multiBot/types/text"
	"multiBot/types/webhook"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
	"github.com/segmentio/kafka-go"
)

func logWaWebhook(c *fiber.Ctx, respData interface{}, reqBody interface{}) {
	logger.LogInfo(c, "RESPONSE", logger.LogReqResp{
		RequestBody:  reqBody,
		ResponseData: respData,
	})
}

func ValidateWaWebhookPost(c *fiber.Ctx) error {
	var waWebhookPayload webhook.WaWebhookPayload
	err := c.BodyParser(&waWebhookPayload)
	if err != nil {
		logWaWebhook(c, nil, waWebhookPayload)
		return c.SendStatus(fiber.StatusOK)
	}

	v := validate.Struct(waWebhookPayload)
	if !v.Validate() {
		logWaWebhook(c, nil, waWebhookPayload)
		return c.SendStatus(fiber.StatusOK)
	}

	c.Locals("waWebhookPayload", waWebhookPayload)
	return c.Next()
}

func WaWebhookPost(c *fiber.Ctx) error {
	waWebhookPayload := c.Locals("waWebhookPayload").(webhook.WaWebhookPayload)

	cfg := config.GetConfig()

	logWaWebhook(c, nil, waWebhookPayload)
	var statusSlice []webhook.ProduceWaMessage
	if waWebhookPayload.Entry != nil && len(waWebhookPayload.Entry) > 0 {
		for _, msg := range waWebhookPayload.Entry {
			if msg.Changes != nil && len(msg.Changes) > 0 {
				for _, change := range msg.Changes {
					changeValue := change.Value
					waAccountPhoneNumberId := changeValue.MetaData.PhoneNumberId
					waAccountInfo := config.GetWhatsappAccountInfoUsingPhoneId(waAccountPhoneNumberId)
					// handler: message status
					if changeValue.Statuses != nil && len(changeValue.Statuses) > 0 {
						for _, status := range changeValue.Statuses {
							if status.Id != "" && status.Status != "" {
								statusSlice = append(statusSlice, webhook.ProduceWaMessage{
									WaId:   status.Id,
									Status: status.Status,
								})
							}
						}
					}
					if !waAccountInfo.OUTGOING_ONLY { // then only we are doing 2 way communication
						// handler: 2 way communication
						if changeValue.Messages != nil && len(changeValue.Messages) > 0 {
							// check if we need to loop on this ?
							firstMessage := changeValue.Messages[0]
							typeOfMsg := firstMessage.Type
							switch typeOfMsg {
							case string(constants.MSG_TYPE_WA_WEBHOOK_TEXT):
								{
									msgText := firstMessage.Text.Body
									return textReplier.ReplyToText(c, msgText, firstMessage.From, waAccountInfo)
								}
							default:
								{
									logger.LogInfo(c, "waWebhook "+"client "+waAccountInfo.NAME+": "+"Received msg of another type", nil)
									textSender.SendTextMessage(
										text.TextDataStruct{
											PreviewUrl: false,
											Body:       "Hi, I think you have given wrong input",
										},
										firstMessage.From,                  // to this number
										string(constants.MSG_PURPOSE_CHAT), // purpose is chat
										constants.SystemClientIdentifier,   // client is SYSTEM i.e. who is sending this message
										waAccountInfo.IDENTIFIER,           // to this wa account
									)
									return c.SendStatus(fiber.StatusOK)
								}
							}
						}
					} else {
						logger.LogInfo(c, "waWebhook "+"client "+waAccountInfo.NAME+": "+"Received msg from user on account where only outgoing can be sent", nil)
						return c.SendStatus(fiber.StatusOK)
					}
				}
			}
		}
	} else {
		return c.SendStatus(fiber.StatusOK)
	}
	// Msg status update logic, in case of 2 way statusSlice is nil so this will not be triggered
	if len(statusSlice) > 0 {
		for _, status := range statusSlice {
			logger.LogInfo(
				c,
				"waWebhook Whatsapp status received at webhook, msgId: "+status.WaId+", status: "+status.Status,
				status,
			)
			if !cfg.WEBHOOK_FALLBACK {
				// TODO: Check for go routines and buffered channel implementation after checking performance
				queue.MessageChannel <- kafka.Message{
					Value: []byte(status.Status),
					Key:   []byte(status.WaId),
				}
			} else {
				service.UpdateWebhookMsgStatus(status.WaId, status.Status)
			}
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
