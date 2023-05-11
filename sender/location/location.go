package location

import (
	"multiBot/config"
	"multiBot/constants"
	"multiBot/models"
	"multiBot/types/location"
	"multiBot/utils"
	"time"
)

func MakeLocationMessage(
	locationBody location.LocationDataStruct,
	to string,
	purpose string,
	clientIdentifier string,
	waAccountIdentifier string,
) (*models.Message, error) {
	var statusHistory []models.StatusHistory
	message, err := utils.ConvertBodyToBsonM(locationBody)
	if err != nil {
		return nil, err
	}
	waAccountInfo := config.GetWhatsappAccountInfo(waAccountIdentifier)
	locationMessage := models.GetNewMessageStructForSending(
		*message,
		clientIdentifier,
		waAccountIdentifier,
		purpose,
		to,
		waAccountInfo.NUMBER,
		constants.MSG_STATUS_CREATED,
		constants.MSG_TYPE_LOCATION,
		constants.MSG_TYPE_LOCATION,
		append(statusHistory, models.StatusHistory{
			Status:    constants.MSG_STATUS_CREATED,
			UpdatedAt: time.Now().UnixMilli(),
		}),
	)
	return &locationMessage, nil
}

func SendLocationMessage(
	locationBody location.LocationDataStruct,
	to string,
	purpose string,
	clientIdentifier string,
	waAccountIdentifier string,
) (*string, error) {
	var messageModel models.Message
	locationMessage, err := MakeLocationMessage(
		locationBody,
		to,
		purpose,
		clientIdentifier,
		waAccountIdentifier,
	)
	if err != nil {
		return nil, err
	}
	objId, err := messageModel.InsertOneMessage(*locationMessage)
	if err != nil {
		return nil, err
	}
	return objId, nil
}
