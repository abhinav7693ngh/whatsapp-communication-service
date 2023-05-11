package documentUrl

import (
	"multiBot/config"
	"multiBot/constants"
	"multiBot/models"
	"multiBot/types/documentUrl"
	"multiBot/utils"
	"time"
)

func MakeDocumentUrlMessage(
	documentUrlBody documentUrl.DocumentUrlDataStruct,
	to string,
	purpose string,
	clientIdentifier string,
	waAccountIdentifier string,
) (*models.Message, error) {
	var statusHistory []models.StatusHistory
	message, err := utils.ConvertBodyToBsonM(documentUrlBody)
	if err != nil {
		return nil, err
	}
	waAccountInfo := config.GetWhatsappAccountInfo(waAccountIdentifier)
	documentUrlMessage := models.GetNewMessageStructForSending(
		*message,
		clientIdentifier,
		waAccountIdentifier,
		purpose,
		to,
		waAccountInfo.NUMBER,
		constants.MSG_STATUS_CREATED,
		constants.MSG_TYPE_DOCUMENT,
		constants.MsgType(constants.MSG_TYPE_DOCUMENT_URL),
		append(statusHistory, models.StatusHistory{
			Status:    constants.MSG_STATUS_CREATED,
			UpdatedAt: time.Now().UnixMilli(),
		}),
	)
	return &documentUrlMessage, nil
}

func SendDocumentUrlMessage(
	documentUrlBody documentUrl.DocumentUrlDataStruct,
	to string,
	purpose string,
	clientIdentifier string,
	waAccountIdentifier string,
) (*string, error) {
	var messageModel models.Message
	documentUrlMessage, err := MakeDocumentUrlMessage(
		documentUrlBody,
		to,
		purpose,
		clientIdentifier,
		waAccountIdentifier,
	)
	if err != nil {
		return nil, err
	}
	objId, err := messageModel.InsertOneMessage(*documentUrlMessage)
	if err != nil {
		return nil, err
	}
	return objId, nil
}
