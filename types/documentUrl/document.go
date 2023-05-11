package documentUrl

import (
	"multiBot/constants"
)

type DocumentUrlDataStruct struct {
	Link    string `json:"link" validate:"required" message:"link is required field"`
	Caption string `json:"caption" validate:"required" message:"caption is required field"`
}

type DocumentUrlPayloadWhatsapp struct {
	MessagingProduct string                 `json:"messaging_product"`
	RecipientType    string                 `json:"recipient_type"`
	To               string                 `json:"to"`
	Type             string                 `json:"type"`
	Document         map[string]interface{} `json:"document"`
}

type DocumentUrlPayloadClient struct {
	To      string                `json:"to" validate:"required" message:"to is required field"`
	Type    string                `json:"type" validate:"required|DocumentUrlTypeValidator" message:"type is required field and should be 'documentUrl' only"`
	Purpose string                `json:"purpose" validate:"required|PurposeValidator" message:"purpose is required and should be either transactional or marketing"`
	Body    DocumentUrlDataStruct `json:"body" validate:"required" message:"body is required field"`
}

func (f DocumentUrlPayloadClient) PurposeValidator(val string) bool {
	return (val == string(constants.MSG_PURPOSE_MARKETING) || val == string(constants.MSG_PURPOSE_TRANSACTIONAL))
}

func (f DocumentUrlPayloadClient) DocumentUrlTypeValidator(val string) bool {
	return val == string(constants.MSG_TYPE_DOCUMENT_URL)
}
