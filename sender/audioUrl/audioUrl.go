package audioUrl

import (
	"multiBot/config"
	"multiBot/constants"
	"multiBot/models"
	"multiBot/types/audioUrl"
	"multiBot/utils"
	"time"
)

func MakeAudioUrlMessage(
	audioUrlBody audioUrl.AudioUrlDataStruct,
	to string,
	purpose string,
	clientIdentifier string,
	waAccountIdentifier string,
) (*models.Message, error) {
	var statusHistory []models.StatusHistory
	message, err := utils.ConvertBodyToBsonM(audioUrlBody)
	if err != nil {
		return nil, err
	}
	waAccountInfo := config.GetWhatsappAccountInfo(waAccountIdentifier)
	audioUrlMessage := models.GetNewMessageStructForSending(
		*message,
		clientIdentifier,
		waAccountIdentifier,
		purpose,
		to,
		waAccountInfo.NUMBER,
		constants.MSG_STATUS_CREATED,
		constants.MSG_TYPE_AUDIO,
		constants.MsgType(constants.MSG_TYPE_AUDIO_URL),
		append(statusHistory, models.StatusHistory{
			Status:    constants.MSG_STATUS_CREATED,
			UpdatedAt: time.Now().UnixMilli(),
		}),
	)
	return &audioUrlMessage, nil
}

func SendAudioUrlMessage(
	audioUrlBody audioUrl.AudioUrlDataStruct,
	to string,
	purpose string,
	clientIdentifier string,
	waAccountIdentifier string,
) (*string, error) {
	var messageModel models.Message
	audioUrlMessage, err := MakeAudioUrlMessage(
		audioUrlBody,
		to,
		purpose,
		clientIdentifier,
		waAccountIdentifier,
	)
	if err != nil {
		return nil, err
	}
	objId, err := messageModel.InsertOneMessage(*audioUrlMessage)
	if err != nil {
		return nil, err
	}
	return objId, nil
}
