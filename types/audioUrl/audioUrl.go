package audioUrl

import (
	"multiBot/constants"
)

type AudioUrlDataStruct struct {
	Link string `json:"link" validate:"required" message:"link is required field"`
}

type AudioUrlPayloadWhatsapp struct {
	MessagingProduct string                 `json:"messaging_product"`
	RecipientType    string                 `json:"recipient_type"`
	To               string                 `json:"to"`
	Type             string                 `json:"type"`
	Audio            map[string]interface{} `json:"audio"`
}

type AudioUrlPayloadClient struct {
	To      string             `json:"to" validate:"required" message:"to is required field"`
	Type    string             `json:"type" validate:"required|AudioUrlTypeValidator" message:"type is required field and should be 'audioUrl' only"`
	Purpose string             `json:"purpose" validate:"required|PurposeValidator" message:"purpose is required and should be either transactional or marketing"`
	Body    AudioUrlDataStruct `json:"body" validate:"required" message:"body is required field"`
}

func (f AudioUrlPayloadClient) PurposeValidator(val string) bool {
	return (val == string(constants.MSG_PURPOSE_MARKETING) || val == string(constants.MSG_PURPOSE_TRANSACTIONAL))
}

func (f AudioUrlPayloadClient) AudioUrlTypeValidator(val string) bool {
	return val == string(constants.MSG_TYPE_AUDIO_URL)
}
