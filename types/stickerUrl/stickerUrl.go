package stickerUrl

import (
	"multiBot/constants"
)

type StickerUrlDataStruct struct {
	Link string `json:"link" validate:"required" message:"link is required field"`
}

type StickerUrlPayloadWhatsapp struct {
	MessagingProduct string                 `json:"messaging_product"`
	RecipientType    string                 `json:"recipient_type"`
	To               string                 `json:"to"`
	Type             string                 `json:"type"`
	Sticker          map[string]interface{} `json:"sticker"`
}

type StickerUrlPayloadClient struct {
	To      string               `json:"to" validate:"required" message:"to is required field"`
	Type    string               `json:"type" validate:"required|StickerUrlTypeValidator" message:"type is required field and should be 'stickerUrl' only"`
	Purpose string               `json:"purpose" validate:"required|PurposeValidator" message:"purpose is required and should be either transactional or marketing"`
	Body    StickerUrlDataStruct `json:"body" validate:"required" message:"body is required field"`
}

func (f StickerUrlPayloadClient) PurposeValidator(val string) bool {
	return (val == string(constants.MSG_PURPOSE_MARKETING) || val == string(constants.MSG_PURPOSE_TRANSACTIONAL))
}

func (f StickerUrlPayloadClient) StickerUrlTypeValidator(val string) bool {
	return val == string(constants.MSG_TYPE_STICKER_URL)
}
