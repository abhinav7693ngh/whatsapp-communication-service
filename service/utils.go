package service

import (
	"encoding/json"
	"multiBot/constants"
	"multiBot/logger"
	"multiBot/models"
	"multiBot/types/audioUrl"
	"multiBot/types/documentUrl"
	"multiBot/types/imageUrl"
	"multiBot/types/list"
	"multiBot/types/location"
	"multiBot/types/replyButton"
	"multiBot/types/stickerUrl"
	"multiBot/types/template"
	"multiBot/types/text"
	"multiBot/types/videoUrl"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PrepareTextTypeMessage(message models.Message) (*[]byte, error) {
	var textWhatsapp text.TextPayloadWhatsapp

	textWhatsapp.Text = message.Body
	textWhatsapp.MessagingProduct = "whatsapp"
	textWhatsapp.RecipientType = "individual"
	textWhatsapp.To = message.Receiver.Contact
	textWhatsapp.Type = string(message.MsgType)

	whatsappJSON, err := json.Marshal(textWhatsapp)
	if err != nil {
		logger.LogError(
			nil,
			"SendWhatsappMessage GoRoutine: error in json marshaling in prepare text type message, id: "+message.Id.Hex()+", error: "+err.Error(),
			nil,
		)
		return nil, err
	}

	return &whatsappJSON, nil
}

func PrepareImageUrlTypeMessage(message models.Message) (*[]byte, error) {
	var imageUrlWhatsapp imageUrl.ImageUrlPayloadWhatsapp

	imageUrlWhatsapp.Image = message.Body
	imageUrlWhatsapp.MessagingProduct = "whatsapp"
	imageUrlWhatsapp.RecipientType = "individual"
	imageUrlWhatsapp.To = message.Receiver.Contact
	imageUrlWhatsapp.Type = string(message.MsgType)

	whatsappJSON, err := json.Marshal(imageUrlWhatsapp)
	if err != nil {
		logger.LogError(
			nil,
			"SendWhatsappMessage GoRoutine: error in json marshaling in prepare image url type message, id: "+message.Id.Hex()+", error: "+err.Error(),
			nil,
		)
		return nil, err
	}

	return &whatsappJSON, nil
}

func PrepareTemplateTypeMessage(message models.Message) (*[]byte, error) {
	var templateWhatsapp template.TemplatePayloadWhatsapp

	templateWhatsapp.Template = message.Body
	templateWhatsapp.MessagingProduct = "whatsapp"
	templateWhatsapp.RecipientType = "individual"
	templateWhatsapp.To = message.Receiver.Contact
	templateWhatsapp.Type = string(message.MsgType)

	whatsappJSON, err := json.Marshal(templateWhatsapp)
	if err != nil {
		logger.LogError(
			nil,
			"SendWhatsappMessage GoRoutine: error in json marshaling in prepare template type message, id: "+message.Id.Hex()+", error: "+err.Error(),
			nil,
		)
		return nil, err
	}

	return &whatsappJSON, nil
}

func PrepareAudioUrlTypeMessage(message models.Message) (*[]byte, error) {
	var audioUrlWhatsapp audioUrl.AudioUrlPayloadWhatsapp

	audioUrlWhatsapp.Audio = message.Body
	audioUrlWhatsapp.MessagingProduct = "whatsapp"
	audioUrlWhatsapp.RecipientType = "individual"
	audioUrlWhatsapp.To = message.Receiver.Contact
	audioUrlWhatsapp.Type = string(message.MsgType)

	whatsappJSON, err := json.Marshal(audioUrlWhatsapp)
	if err != nil {
		logger.LogError(
			nil,
			"SendWhatsappMessage GoRoutine: error in json marshaling in prepare audio url type message, id: "+message.Id.Hex()+", error: "+err.Error(),
			nil,
		)
		return nil, err
	}

	return &whatsappJSON, nil
}

func PrepareDocumentUrlTypeMessage(message models.Message) (*[]byte, error) {
	var documentUrlWhatsapp documentUrl.DocumentUrlPayloadWhatsapp

	documentUrlWhatsapp.Document = message.Body
	documentUrlWhatsapp.MessagingProduct = "whatsapp"
	documentUrlWhatsapp.RecipientType = "individual"
	documentUrlWhatsapp.To = message.Receiver.Contact
	documentUrlWhatsapp.Type = string(message.MsgType)

	whatsappJSON, err := json.Marshal(documentUrlWhatsapp)
	if err != nil {
		logger.LogError(
			nil,
			"SendWhatsappMessage GoRoutine: error in json marshaling in prepare document url type message, id: "+message.Id.Hex()+", error: "+err.Error(),
			nil,
		)
		return nil, err
	}

	return &whatsappJSON, nil
}

func PrepareStickerUrlTypeMessage(message models.Message) (*[]byte, error) {
	var stickerUrlWhatsapp stickerUrl.StickerUrlPayloadWhatsapp

	stickerUrlWhatsapp.Sticker = message.Body
	stickerUrlWhatsapp.MessagingProduct = "whatsapp"
	stickerUrlWhatsapp.RecipientType = "individual"
	stickerUrlWhatsapp.To = message.Receiver.Contact
	stickerUrlWhatsapp.Type = string(message.MsgType)

	whatsappJSON, err := json.Marshal(stickerUrlWhatsapp)
	if err != nil {
		logger.LogError(
			nil,
			"SendWhatsappMessage GoRoutine: error in json marshaling in prepare sticker url type message, id: "+message.Id.Hex()+", error: "+err.Error(),
			nil,
		)
		return nil, err
	}

	return &whatsappJSON, nil
}

func PrepareVideoUrlTypeMessage(message models.Message) (*[]byte, error) {
	var videoUrlWhatsapp videoUrl.VideoUrlPayloadWhatsapp

	videoUrlWhatsapp.Video = message.Body
	videoUrlWhatsapp.MessagingProduct = "whatsapp"
	videoUrlWhatsapp.RecipientType = "individual"
	videoUrlWhatsapp.To = message.Receiver.Contact
	videoUrlWhatsapp.Type = string(message.MsgType)

	whatsappJSON, err := json.Marshal(videoUrlWhatsapp)
	if err != nil {
		logger.LogError(
			nil,
			"SendWhatsappMessage GoRoutine: error in json marshaling in prepare video url type message, id: "+message.Id.Hex()+", error: "+err.Error(),
			nil,
		)
		return nil, err
	}

	return &whatsappJSON, nil
}

func PrepareLocationTypeMessage(message models.Message) (*[]byte, error) {
	var locationWhatsapp location.LocationPayloadWhatsapp

	locationWhatsapp.Location = message.Body
	locationWhatsapp.MessagingProduct = "whatsapp"
	locationWhatsapp.RecipientType = "individual"
	locationWhatsapp.To = message.Receiver.Contact
	locationWhatsapp.Type = string(message.MsgType)

	whatsappJSON, err := json.Marshal(locationWhatsapp)
	if err != nil {
		logger.LogError(
			nil,
			"SendWhatsappMessage GoRoutine: error in json marshaling in prepare location type message, id: "+message.Id.Hex()+", error: "+err.Error(),
			nil,
		)
		return nil, err
	}

	return &whatsappJSON, nil
}

func PrepareListTypeMessage(message models.Message) (*[]byte, error) {
	var listWhatsapp list.ListPayloadWhatsapp

	listWhatsapp.Interactive = message.Body
	listWhatsapp.MessagingProduct = "whatsapp"
	listWhatsapp.RecipientType = "individual"
	listWhatsapp.To = message.Receiver.Contact
	listWhatsapp.Type = string(message.MsgType)

	whatsappJSON, err := json.Marshal(listWhatsapp)
	if err != nil {
		logger.LogError(
			nil,
			"SendWhatsappMessage GoRoutine: error in json marshaling in prepare list type message, id: "+message.Id.Hex()+", error: "+err.Error(),
			nil,
		)
		return nil, err
	}

	return &whatsappJSON, nil
}

func PrepareReplyButtonTypeMessage(message models.Message) (*[]byte, error) {
	var buttonWhatsapp replyButton.ReplyButtonPayloadWhatsapp

	buttonWhatsapp.Interactive = message.Body
	buttonWhatsapp.MessagingProduct = "whatsapp"
	buttonWhatsapp.RecipientType = "individual"
	buttonWhatsapp.To = message.Receiver.Contact
	buttonWhatsapp.Type = string(message.MsgType)

	whatsappJSON, err := json.Marshal(buttonWhatsapp)
	if err != nil {
		logger.LogError(
			nil,
			"SendWhatsappMessage GoRoutine: error in json marshaling in prepare reply button type message, id: "+message.Id.Hex()+", error: "+err.Error(),
			nil,
		)
		return nil, err
	}

	return &whatsappJSON, nil
}

func StatusUpdateWhatsappPushFailed(msg [2]string, filter primitive.M, message models.Message) {
	possible := message.PossibleToUpdateStatus(constants.MSG_STATUS_WA_PUSH_FAILED, "SendWhatsappMessage GoRoutine")
	if !possible {
		return
	}
	currentTime := time.Now().UnixMilli()
	update := bson.M{
		"$set": bson.M{
			"status":    constants.MSG_STATUS_WA_PUSH_FAILED,
			"updatedAt": currentTime,
		},
		"$push": bson.M{
			"statusHistory": bson.M{
				"status":    constants.MSG_STATUS_WA_PUSH_FAILED,
				"updatedAt": currentTime,
			},
		},
		"$inc": bson.M{
			"retryCount": 1,
		},
	}
	var messageModel models.Message
	updatedNo, err := messageModel.UpdateOneMessage(filter, update)
	if err != nil {
		// error case
		logger.LogError(
			nil,
			msg[0]+message.Id.Hex()+", error: "+err.Error(),
			nil,
		)
	} else if updatedNo != nil {
		if *updatedNo > 0 {
			// updated case
			logger.LogInfo(
				nil,
				msg[1]+message.Id.Hex(),
				nil,
			)
		}
	}
}
