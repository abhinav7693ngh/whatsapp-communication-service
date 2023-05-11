package videoUrl

import (
	"multiBot/config"
	"multiBot/constants"
	"multiBot/models"
	"multiBot/types/videoUrl"
	"multiBot/utils"
	"time"
)

func MakeVideoUrlMessage(
	videoUrlBody videoUrl.VideoUrlDataStruct,
	to string,
	purpose string,
	clientIdentifier string,
	waAccountIdentifier string,
) (*models.Message, error) {
	var statusHistory []models.StatusHistory
	message, err := utils.ConvertBodyToBsonM(videoUrlBody)
	if err != nil {
		return nil, err
	}
	waAccountInfo := config.GetWhatsappAccountInfo(waAccountIdentifier)
	videoUrlMessage := models.GetNewMessageStructForSending(
		*message,
		clientIdentifier,
		waAccountIdentifier,
		purpose,
		to,
		waAccountInfo.NUMBER,
		constants.MSG_STATUS_CREATED,
		constants.MSG_TYPE_VIDEO,
		constants.MsgType(constants.MSG_TYPE_VIDEO_URL),
		append(statusHistory, models.StatusHistory{
			Status:    constants.MSG_STATUS_CREATED,
			UpdatedAt: time.Now().UnixMilli(),
		}),
	)
	return &videoUrlMessage, nil
}

func SendVideoUrlMessage(
	videoUrlBody videoUrl.VideoUrlDataStruct,
	to string,
	purpose string,
	clientIdentifier string,
	waAccountIdentifier string,
) (*string, error) {
	var messageModel models.Message
	videoUrlMessage, err := MakeVideoUrlMessage(
		videoUrlBody,
		to,
		purpose,
		clientIdentifier,
		waAccountIdentifier,
	)
	if err != nil {
		return nil, err
	}
	objId, err := messageModel.InsertOneMessage(*videoUrlMessage)
	if err != nil {
		return nil, err
	}
	return objId, nil
}
