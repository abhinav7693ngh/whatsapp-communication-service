package videoUrl

import (
	"multiBot/constants"
)

type VideoUrlDataStruct struct {
	Link    string `json:"link" validate:"required" message:"link is required field"`
	Caption string `json:"caption" validate:"required" message:"caption is required field"`
}

type VideoUrlPayloadWhatsapp struct {
	MessagingProduct string                 `json:"messaging_product"`
	RecipientType    string                 `json:"recipient_type"`
	To               string                 `json:"to"`
	Type             string                 `json:"type"`
	Video            map[string]interface{} `json:"video"`
}

type VideoUrlPayloadClient struct {
	To      string             `json:"to" validate:"required" message:"to is required field"`
	Type    string             `json:"type" validate:"required|VideoUrlTypeValidator" message:"type is required field and should be 'videoUrl' only"`
	Purpose string             `json:"purpose" validate:"required|PurposeValidator" message:"purpose is required and should be either transactional or marketing"`
	Body    VideoUrlDataStruct `json:"body" validate:"required" message:"body is required field"`
}

func (f VideoUrlPayloadClient) PurposeValidator(val string) bool {
	return (val == string(constants.MSG_PURPOSE_MARKETING) || val == string(constants.MSG_PURPOSE_TRANSACTIONAL))
}

func (f VideoUrlPayloadClient) VideoUrlTypeValidator(val string) bool {
	return val == string(constants.MSG_TYPE_VIDEO_URL)
}
