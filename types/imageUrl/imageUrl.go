package imageUrl

import (
	"multiBot/constants"
)

type ImageUrlDataStruct struct {
	Link string `json:"link" validate:"required" message:"link is required field"`
}

type ImageUrlPayloadWhatsapp struct {
	MessagingProduct string                 `json:"messaging_product"`
	RecipientType    string                 `json:"recipient_type"`
	To               string                 `json:"to"`
	Type             string                 `json:"type"`
	Image            map[string]interface{} `json:"image"`
}

type ImageUrlPayloadClient struct {
	To      string             `json:"to" validate:"required" message:"to is required field"`
	Type    string             `json:"type" validate:"required|ImageUrlTypeValidator" message:"type is required field and should be 'imageUrl' only"`
	Purpose string             `json:"purpose" validate:"required|PurposeValidator" message:"purpose is required and should be either transactional or marketing"`
	Body    ImageUrlDataStruct `json:"body" validate:"required" message:"body is required field"`
}

func (f ImageUrlPayloadClient) PurposeValidator(val string) bool {
	return (val == string(constants.MSG_PURPOSE_MARKETING) || val == string(constants.MSG_PURPOSE_TRANSACTIONAL))
}

func (f ImageUrlPayloadClient) ImageUrlTypeValidator(val string) bool {
	return val == string(constants.MSG_TYPE_IMAGE_URL)
}
