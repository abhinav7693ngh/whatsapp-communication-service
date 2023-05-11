package stickerUrl

import (
	"multiBot/config"
	"multiBot/constants"
	"multiBot/models"
	"multiBot/types/stickerUrl"
	"multiBot/utils"
	"time"
)

func MakeStickerUrlMessage(
	stickerUrlBody stickerUrl.StickerUrlDataStruct,
	to string,
	purpose string,
	clientIdentifier string,
	waAccountIdentifier string,
) (*models.Message, error) {
	var statusHistory []models.StatusHistory
	message, err := utils.ConvertBodyToBsonM(stickerUrlBody)
	if err != nil {
		return nil, err
	}
	waAccountInfo := config.GetWhatsappAccountInfo(waAccountIdentifier)
	stickerUrlMessage := models.GetNewMessageStructForSending(
		*message,
		clientIdentifier,
		waAccountIdentifier,
		purpose,
		to,
		waAccountInfo.NUMBER,
		constants.MSG_STATUS_CREATED,
		constants.MSG_TYPE_STICKER,
		constants.MsgType(constants.MSG_TYPE_STICKER_URL),
		append(statusHistory, models.StatusHistory{
			Status:    constants.MSG_STATUS_CREATED,
			UpdatedAt: time.Now().UnixMilli(),
		}),
	)
	return &stickerUrlMessage, nil
}

func SendStickerUrlMessage(
	stickerUrlBody stickerUrl.StickerUrlDataStruct,
	to string,
	purpose string,
	clientIdentifier string,
	waAccountIdentifier string,
) (*string, error) {
	var messageModel models.Message
	stickerUrlMessage, err := MakeStickerUrlMessage(
		stickerUrlBody,
		to,
		purpose,
		clientIdentifier,
		waAccountIdentifier,
	)
	if err != nil {
		return nil, err
	}
	objId, err := messageModel.InsertOneMessage(*stickerUrlMessage)
	if err != nil {
		return nil, err
	}
	return objId, nil
}
