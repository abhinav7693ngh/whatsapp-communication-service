package template

import (
	"multiBot/config"
	"multiBot/constants"
	"multiBot/models"
	templateParser "multiBot/parser/template"
	"multiBot/types/template"
	"multiBot/utils"
	"time"
)

func MakeTemplateMessage(
	templateBody template.TemplateDataStruct,
	to string,
	purpose string,
	clientIdentifier string,
	waAccountIdentifier string,
) (*models.Message, error) {
	var statusHistory []models.StatusHistory
	msgBody := templateParser.ParseTemplate(templateBody)
	message, err := utils.ConvertBodyToBsonM(msgBody)
	if err != nil {
		return nil, err
	}
	waAccountInfo := config.GetWhatsappAccountInfo(waAccountIdentifier)
	templateMessage := models.GetNewMessageStructForSending(
		*message,
		clientIdentifier,
		waAccountIdentifier,
		purpose,
		to,
		waAccountInfo.NUMBER,
		constants.MSG_STATUS_CREATED,
		constants.MSG_TYPE_TEMPLATE,
		constants.MSG_TYPE_TEMPLATE,
		append(statusHistory, models.StatusHistory{
			Status:    constants.MSG_STATUS_CREATED,
			UpdatedAt: time.Now().UnixMilli(),
		}),
	)
	return &templateMessage, nil
}

func SendTemplateMessage(
	templateBody template.TemplateDataStruct,
	to string,
	purpose string,
	clientIdentifier string,
	waAccountIdentifier string,
) (*string, error) {
	var messageModel models.Message
	templateMessage, err := MakeTemplateMessage(
		templateBody,
		to,
		purpose,
		clientIdentifier,
		waAccountIdentifier,
	)
	if err != nil {
		return nil, err
	}
	objId, err := messageModel.InsertOneMessage(*templateMessage)
	if err != nil {
		return nil, err
	}
	return objId, nil
}
