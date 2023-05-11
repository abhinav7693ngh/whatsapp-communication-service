package webhook

type WaStatusObj struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

type TextMessageTwoWay struct {
	Body string `json:"body"`
}

type WaMessagesTwoWay struct {
	From string            `json:"from"`
	Type string            `json:"type"`
	Text TextMessageTwoWay `json:"text"`
}

type WaMetaData struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberId      string `json:"phone_number_id"`
}

type WaValueData struct {
	MetaData WaMetaData         `json:"metadata"`
	Statuses []WaStatusObj      `json:"statuses"`
	Messages []WaMessagesTwoWay `json:"messages"`
}

type WaValue struct {
	Value WaValueData `json:"value" validate:"required"`
}

type WaChanges struct {
	Changes []WaValue `json:"changes" validate:"required"`
}

type WaWebhookPayload struct {
	Entry []WaChanges `json:"entry" validate:"required"`
}
