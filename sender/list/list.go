package list

import (
	"multiBot/config"
	"multiBot/constants"
	"multiBot/models"
	listParser "multiBot/parser/list"
	"multiBot/types/list"
	"multiBot/utils"
	"time"
)

func MakeListMessage(
	listBody list.ListDataStruct,
	to string,
	purpose string,
	clientIdentifier string,
	waAccountIdentifier string,
) (*models.Message, error) {
	var statusHistory []models.StatusHistory
	msgBody := listParser.ParseList(listBody)
	message, err := utils.ConvertBodyToBsonM(msgBody)
	if err != nil {
		return nil, err
	}
	waAccountInfo := config.GetWhatsappAccountInfo(waAccountIdentifier)
	listMessage := models.GetNewMessageStructForSending(
		*message,
		clientIdentifier,
		waAccountIdentifier,
		purpose,
		to,
		waAccountInfo.NUMBER,
		constants.MSG_STATUS_CREATED,
		constants.MSG_TYPE_INTERACTIVE,
		constants.MsgType(constants.MSG_TYPE_LIST_INTERNAL),
		append(statusHistory, models.StatusHistory{
			Status:    constants.MSG_STATUS_CREATED,
			UpdatedAt: time.Now().UnixMilli(),
		}),
	)
	return &listMessage, nil
}

func SendListMessage(
	listBody list.ListDataStruct,
	to string,
	purpose string,
	clientIdentifier string,
	waAccountIdentifier string,
) (*string, error) {
	var messageModel models.Message
	listMessage, err := MakeListMessage(
		listBody,
		to,
		purpose,
		clientIdentifier,
		waAccountIdentifier,
	)
	if err != nil {
		return nil, err
	}
	objId, err := messageModel.InsertOneMessage(*listMessage)
	if err != nil {
		return nil, err
	}
	return objId, nil
}
