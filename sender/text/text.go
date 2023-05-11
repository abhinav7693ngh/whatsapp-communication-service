package text

import (
	"multiBot/config"
	"multiBot/constants"
	"multiBot/models"
	"multiBot/types/text"
	"multiBot/utils"
	"time"
)

func MakeTextMessage(
	textBody text.TextDataStruct,
	to string,
	purpose string,
	clientIdentifier string,
	waAccountIdentifier string,
) (*models.Message, error) {
	var statusHistory []models.StatusHistory
	message, err := utils.ConvertBodyToBsonM(textBody)
	if err != nil {
		return nil, err
	}
	waAccountInfo := config.GetWhatsappAccountInfo(waAccountIdentifier)
	textMessage := models.GetNewMessageStructForSending(
		*message,
		clientIdentifier,
		waAccountIdentifier,
		purpose,
		to,
		waAccountInfo.NUMBER,
		constants.MSG_STATUS_CREATED,
		constants.MSG_TYPE_TEXT,
		constants.MSG_TYPE_TEXT,
		append(statusHistory, models.StatusHistory{
			Status:    constants.MSG_STATUS_CREATED,
			UpdatedAt: time.Now().UnixMilli(),
		}),
	)
	return &textMessage, nil
}

func SendTextMessage(
	textBody text.TextDataStruct,
	to string,
	purpose string,
	clientIdentifier string,
	waAccountIdentifier string,
) (*string, error) {
	var messageModel models.Message
	textMessage, err := MakeTextMessage(
		textBody,
		to,
		purpose,
		clientIdentifier,
		waAccountIdentifier,
	)
	if err != nil {
		return nil, err
	}
	objId, err := messageModel.InsertOneMessage(*textMessage)
	if err != nil {
		return nil, err
	}
	return objId, nil
}
