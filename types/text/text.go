package text

import (
	"multiBot/constants"
)

type TextDataStruct struct {
	PreviewUrl bool   `json:"preview_url" bson:"preview_url"`
	Body       string `json:"body" validate:"required" bson:"body"`
}

type TextPayloadWhatsapp struct {
	MessagingProduct string                 `json:"messaging_product"`
	RecipientType    string                 `json:"recipient_type"`
	To               string                 `json:"to"`
	Type             string                 `json:"type"`
	Text             map[string]interface{} `json:"text"`
}

type TextPayloadClient struct {
	To      string         `json:"to" validate:"required" message:"to is required field"`
	Type    string         `json:"type" validate:"required|TextTypeValidator" message:"type is required field and should be of 'text' only"`
	Purpose string         `json:"purpose" validate:"required|PurposeValidator" message:"purpose is required and should be either transactional or marketing"`
	Body    TextDataStruct `json:"body" validate:"required" message:"body is required field"`
}

func (f TextPayloadClient) PurposeValidator(val string) bool {
	return (val == string(constants.MSG_PURPOSE_MARKETING) || val == string(constants.MSG_PURPOSE_TRANSACTIONAL))
}

func (f TextPayloadClient) TextTypeValidator(val string) bool {
	return val == string(constants.MSG_TYPE_TEXT)
}
