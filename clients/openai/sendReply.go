package openai

import (
	"context"
	"multiBot/config"
	"multiBot/constants"
	"multiBot/logger"
	"multiBot/models"
	textSender "multiBot/sender/text"
	"multiBot/types/text"

	"github.com/gofiber/fiber/v2"
	openai "github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func reverseArrayOpenAi(arr []openai.ChatCompletionMessage) []openai.ChatCompletionMessage {
	for i := 0; i < len(arr)/2; i++ {
		j := len(arr) - i - 1
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func SendReply(c *fiber.Ctx, msgText string, from string, waAccount *config.WhatsappAccountsConfig) error {

	cfg := config.GetConfig()
	openAiToken := cfg.OPEN_AI.TOKEN
	aiTraining := cfg.OPEN_AI.TRAINING

	// 1. Get last 4 msgs for that user from DB which we have sent
	// 2. make an array of messages
	// 3. send to gpt for chat context
	filter := bson.M{
		"receiver.contact":          from,
		"internalType":              constants.MSG_TYPE_TEXT,
		"whatsappAccountIdentifier": waAccount.IDENTIFIER,
		"msgPurpose":                string(constants.MSG_PURPOSE_CHAT),
	}
	options := options.Find().SetSort(bson.M{"createdAt": -1}).SetLimit(4)
	var messageModel models.Message
	msgs, err := messageModel.FindManyMessages(filter, options)
	if err != nil {
		logger.LogError(c, "waWebhook "+"client "+waAccount.NAME+": "+"Get 4 msgs from DB error: "+err.Error(), nil)
		textSender.SendTextMessage(
			text.TextDataStruct{
				PreviewUrl: false,
				Body:       "I am having trouble answering, Please try again later",
			},
			from,                               // to this number
			string(constants.MSG_PURPOSE_CHAT), // purpose is chat
			constants.SystemClientIdentifier,   // client is SYSTEM i.e. who is sending this message
			waAccount.IDENTIFIER,               // to this wa account
		)
		return c.SendStatus(fiber.StatusOK)
	}

	// openai messages creation for context
	openAiMsgs := []openai.ChatCompletionMessage{}
	for _, msg := range *msgs {
		var msgRole string
		if msg.MsgNetworkType == string(constants.MSG_NETWORK_OUTGOING) {
			msgRole = openai.ChatMessageRoleAssistant
		} else if msg.MsgNetworkType == string(constants.MSG_NETWORK_INCOMING) {
			msgRole = openai.ChatMessageRoleUser
		}
		textMsgBody, _ := msg.Body["body"].(string)
		openAiMsgs = append(openAiMsgs, openai.ChatCompletionMessage{
			Role:    msgRole,
			Content: textMsgBody,
		})
	}

	// reversed msgs to maintain order
	reversedOpenAiMsgs := reverseArrayOpenAi(openAiMsgs)

	// Append with new msgs from user
	reversedOpenAiMsgs = append(reversedOpenAiMsgs, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: msgText,
	})

	// prepend with training data
	reversedOpenAiMsgs = append([]openai.ChatCompletionMessage{{
		Role:    openai.ChatMessageRoleSystem,
		Content: aiTraining,
	}}, reversedOpenAiMsgs...)

	// OpenAI Client
	client := openai.NewClient(openAiToken)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Temperature:      0.45,
			MaxTokens:        512,
			Model:            openai.GPT3Dot5Turbo,
			Messages:         reversedOpenAiMsgs,
			TopP:             1,
			FrequencyPenalty: 0,
			PresencePenalty:  0,
		},
	)
	if err != nil {
		logger.LogError(c, "waWebhook "+"client "+waAccount.NAME+": "+err.Error(), nil)
		textSender.SendTextMessage(
			text.TextDataStruct{
				PreviewUrl: false,
				Body:       "I am having trouble answering, Please try again later",
			},
			from,                               // to this number
			string(constants.MSG_PURPOSE_CHAT), // purpose is chat
			constants.SystemClientIdentifier,   // client is SYSTEM i.e. who is sending this message
			waAccount.IDENTIFIER,               // to this wa account
		)
		return c.SendStatus(fiber.StatusOK)
	}

	if len(resp.Choices) > 0 {
		msgToSend := resp.Choices[0].Message.Content
		textSender.SendTextMessage(
			text.TextDataStruct{
				PreviewUrl: false,
				Body:       msgToSend,
			},
			from,                               // to this number
			string(constants.MSG_PURPOSE_CHAT), // purpose is chat
			constants.SystemClientIdentifier,   // client is SYSTEM i.e. who is sending this message
			waAccount.IDENTIFIER,               // to this wa account
		)
		return c.SendStatus(fiber.StatusOK)
	} else {
		logger.LogError(c, "waWebhook "+"client "+waAccount.NAME+": "+"Did not receive any message", nil)
		textSender.SendTextMessage(
			text.TextDataStruct{
				PreviewUrl: false,
				Body:       "I am having trouble answering, Please try again later",
			},
			from,                               // to this number
			string(constants.MSG_PURPOSE_CHAT), // purpose is chat
			constants.SystemClientIdentifier,   // client is SYSTEM i.e. who is sending this message
			waAccount.IDENTIFIER,               // to this wa account
		)
		return c.SendStatus(fiber.StatusOK)
	}
}
