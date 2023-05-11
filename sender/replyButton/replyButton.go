package replyButton

import (
	"multiBot/config"
	"multiBot/constants"
	"multiBot/models"
	replyButtonParser "multiBot/parser/replyButton"
	"multiBot/types/replyButton"
	"multiBot/utils"
	"time"
)

func MakeReplyButtonMessage(
	replyButtonBody replyButton.ReplyButtonDataStruct,
	to string,
	purpose string,
	clientIdentifier string,
	waAccountIdentifier string,
) (*models.Message, error) {
	var statusHistory []models.StatusHistory
	msgBody := replyButtonParser.ParseReplyButton(replyButtonBody)
	message, err := utils.ConvertBodyToBsonM(msgBody)
	if err != nil {
		return nil, err
	}
	waAccountInfo := config.GetWhatsappAccountInfo(waAccountIdentifier)
	replyButtonMessage := models.GetNewMessageStructForSending(
		*message,
		clientIdentifier,
		waAccountIdentifier,
		purpose,
		to,
		waAccountInfo.NUMBER,
		constants.MSG_STATUS_CREATED,
		constants.MSG_TYPE_INTERACTIVE,
		constants.MsgType(constants.MSG_TYPE_REPLY_BUTTON_INTERNAL),
		append(statusHistory, models.StatusHistory{
			Status:    constants.MSG_STATUS_CREATED,
			UpdatedAt: time.Now().UnixMilli(),
		}),
	)
	return &replyButtonMessage, nil
}

func SendReplyButtonMessage(
	replyButtonBody replyButton.ReplyButtonDataStruct,
	to string,
	purpose string,
	clientIdentifier string,
	waAccountIdentifier string,
) (*string, error) {
	var messageModel models.Message
	replyButtonMessage, err := MakeReplyButtonMessage(
		replyButtonBody,
		to,
		purpose,
		clientIdentifier,
		waAccountIdentifier,
	)
	if err != nil {
		return nil, err
	}
	objId, err := messageModel.InsertOneMessage(*replyButtonMessage)
	if err != nil {
		return nil, err
	}
	return objId, nil
}
