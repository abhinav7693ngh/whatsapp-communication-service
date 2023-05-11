package imageUrl

import (
	"multiBot/config"
	"multiBot/constants"
	"multiBot/models"
	"multiBot/types/imageUrl"
	"multiBot/utils"
	"time"
)

func MakeImageUrlMessage(
	imageUrlBody imageUrl.ImageUrlDataStruct,
	to string,
	purpose string,
	clientIdentifier string,
	waAccountIdentifier string,
) (*models.Message, error) {
	var statusHistory []models.StatusHistory
	message, err := utils.ConvertBodyToBsonM(imageUrlBody)
	if err != nil {
		return nil, err
	}
	waAccountInfo := config.GetWhatsappAccountInfo(waAccountIdentifier)
	imageUrlMessage := models.GetNewMessageStructForSending(
		*message,
		clientIdentifier,
		waAccountIdentifier,
		purpose,
		to,
		waAccountInfo.NUMBER,
		constants.MSG_STATUS_CREATED,
		constants.MSG_TYPE_IMAGE,
		constants.MsgType(constants.MSG_TYPE_IMAGE_URL),
		append(statusHistory, models.StatusHistory{
			Status:    constants.MSG_STATUS_CREATED,
			UpdatedAt: time.Now().UnixMilli(),
		}),
	)
	return &imageUrlMessage, nil
}

func SendImageUrlMessage(
	imageUrlBody imageUrl.ImageUrlDataStruct,
	to string,
	purpose string,
	clientIdentifier string,
	waAccountIdentifier string,
) (*string, error) {
	var messageModel models.Message
	imageUrlMessage, err := MakeImageUrlMessage(
		imageUrlBody,
		to,
		purpose,
		clientIdentifier,
		waAccountIdentifier,
	)
	objId, err := messageModel.InsertOneMessage(*imageUrlMessage)
	if err != nil {
		return nil, err
	}
	return objId, nil
}
