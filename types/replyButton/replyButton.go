package replyButton

import (
	"multiBot/constants"
)

type ReplyButtonDataStruct struct {
	Body    string   `json:"body" validate:"required" message:"body is required field"`
	Buttons []string `json:"buttons" validate:"required" message:"buttons is required field"`
}

type ReplyButtonPayloadWhatsapp struct {
	MessagingProduct string                 `json:"messaging_product"`
	RecipientType    string                 `json:"recipient_type"`
	To               string                 `json:"to"`
	Type             string                 `json:"type"`
	Interactive      map[string]interface{} `json:"interactive"`
}

type ReplyButtonPayloadClient struct {
	To      string                `json:"to" validate:"required" message:"to is required field"`
	Type    string                `json:"type" validate:"required|ReplyButtonTypeValidator" message:"type is required field and should be 'replyButton' only"`
	Purpose string                `json:"purpose" validate:"required|PurposeValidator" message:"purpose is required and should be either transactional or marketing"`
	Body    ReplyButtonDataStruct `json:"body" validate:"required" message:"body is required field"`
}

func (f ReplyButtonPayloadClient) PurposeValidator(val string) bool {
	return (val == string(constants.MSG_PURPOSE_MARKETING) || val == string(constants.MSG_PURPOSE_TRANSACTIONAL))
}

func (f ReplyButtonPayloadClient) ReplyButtonTypeValidator(val string) bool {
	return val == string(constants.MSG_TYPE_REPLY_BUTTON_INTERNAL)
}
