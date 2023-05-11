package location

import (
	"multiBot/constants"
)

type LocationDataStruct struct {
	Latitude  float64 `json:"latitude" validate:"required" message:"latitude is required field"`
	Longitude float64 `json:"longitude" validate:"required" message:"longitude is required field"`
	Name      string  `json:"name" validate:"required" message:"name is required field"`
	Address   string  `json:"address" validate:"required" message:"address is required field"`
}

type LocationPayloadWhatsapp struct {
	MessagingProduct string                 `json:"messaging_product"`
	RecipientType    string                 `json:"recipient_type"`
	To               string                 `json:"to"`
	Type             string                 `json:"type"`
	Location         map[string]interface{} `json:"location"`
}

type LocationPayloadClient struct {
	To      string             `json:"to" validate:"required" message:"to is required field"`
	Type    string             `json:"type" validate:"required|LocationTypeValidator" message:"type is required field and should be 'location' only"`
	Purpose string             `json:"purpose" validate:"required|PurposeValidator" message:"purpose is required and should be either transactional or marketing"`
	Body    LocationDataStruct `json:"body" validate:"required" message:"body is required field"`
}

func (f LocationPayloadClient) PurposeValidator(val string) bool {
	return (val == string(constants.MSG_PURPOSE_MARKETING) || val == string(constants.MSG_PURPOSE_TRANSACTIONAL))
}

func (f LocationPayloadClient) LocationTypeValidator(val string) bool {
	return val == string(constants.MSG_TYPE_LOCATION)
}
