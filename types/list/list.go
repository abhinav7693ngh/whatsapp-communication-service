package list

import (
	"multiBot/constants"
)

type ListSectionRows struct {
	Title       string `json:"title" validate:"required" message:"title is required field"`
	Description string `json:"description"`
}

type ListSection struct {
	Title string            `json:"title" validate:"required" message:"title is required field"`
	Rows  []ListSectionRows `json:"rows" validate:"required" message:"rows is required field"`
}

type ListDataStruct struct {
	Header   string        `json:"header"`
	Body     string        `json:"body" validate:"required" message:"body is required field"`
	Footer   string        `json:"footer"`
	Button   string        `json:"button" validate:"required" message:"button is required field"`
	Sections []ListSection `json:"sections" validate:"required" message:"sections is required field"`
}

type ListPayloadWhatsapp struct {
	MessagingProduct string                 `json:"messaging_product"`
	RecipientType    string                 `json:"recipient_type"`
	To               string                 `json:"to"`
	Type             string                 `json:"type"`
	Interactive      map[string]interface{} `json:"interactive"`
}

type ListPayloadClient struct {
	To      string         `json:"to" validate:"required" message:"to is required field"`
	Type    string         `json:"type" validate:"required|ListTypeValidator" message:"type is required field and should be 'list' only"`
	Purpose string         `json:"purpose" validate:"required|PurposeValidator" message:"purpose is required and should be either transactional or marketing"`
	Body    ListDataStruct `json:"body" validate:"required" message:"body is required field"`
}

func (f ListPayloadClient) PurposeValidator(val string) bool {
	return (val == string(constants.MSG_PURPOSE_MARKETING) || val == string(constants.MSG_PURPOSE_TRANSACTIONAL))
}

func (f ListPayloadClient) ListTypeValidator(val string) bool {
	return val == string(constants.MSG_TYPE_LIST_INTERNAL)
}
