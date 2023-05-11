package template

import (
	"encoding/json"
	"errors"
	"multiBot/constants"
)

type HeaderImageData struct {
	Link *string `json:"link"`
}

type HeaderAudioData struct {
	Link *string `json:"link"`
}

type HeaderVideoData struct {
	Link    *string `json:"link"`
	Caption *string `json:"caption"`
}

type HeaderDocumentData struct {
	Link    *string `json:"link"`
	Caption *string `json:"caption"`
}

type HeaderLocationData struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Name      *string  `json:"name"`
	Address   *string  `json:"address"`
}

type BodyTextData struct {
	Text *string `json:"text"`
}

type ButtonQuickReplyData struct {
	Payload *string `json:"payload"`
}

type ButtonUrlData struct {
	Text *string `json:"text"`
}

type TemplateTypeData struct {
	Type *string     `json:"type"`
	Data interface{} `json:"data"`
}

type TemplateDataStruct struct {
	Name     string             `json:"name" validate:"required" message:"name is required field"`
	Language string             `json:"language" validate:"required" message:"language is required field"`
	Header   *TemplateTypeData  `json:"header"`
	Body     []TemplateTypeData `json:"body"`
	Button   []TemplateTypeData `json:"button"`
}

func (t *TemplateDataStruct) UnmarshalJSON(data []byte) error {

	type TemplateTypeDataRaw struct {
		Type *string          `json:"type"`
		Data *json.RawMessage `json:"data"`
	}

	var rawData struct {
		Name     string                `json:"name"`
		Language string                `json:"language"`
		Header   *TemplateTypeDataRaw  `json:"header"`
		Body     []TemplateTypeDataRaw `json:"body"`
		Button   []TemplateTypeDataRaw `json:"button"`
	}
	err := json.Unmarshal(data, &rawData)
	if err != nil {
		return err
	}
	t.Name = rawData.Name
	t.Language = rawData.Language

	// Unmarshal header
	// Header is used for image, location, audio, video, document i.e. for media ( non-text )
	if rawData.Header != nil {
		if rawData.Header.Type == nil {
			return errors.New("header type is required")
		}
		if rawData.Header.Data == nil {
			return errors.New("header data is required")
		}
		switch *rawData.Header.Type {
		case string(constants.TEMPLATE_MSG_TYPE_HEADER_IMAGE): // Header image
			{
				var headerImageData HeaderImageData
				err := json.Unmarshal(*rawData.Header.Data, &headerImageData)
				if err != nil {
					return err
				}
				if headerImageData.Link == nil {
					return errors.New("header image link is required")
				}
				t.Header = &TemplateTypeData{
					Type: rawData.Header.Type,
					Data: headerImageData,
				}
			}
		case string(constants.TEMPLATE_MSG_TYPE_HEADER_LOCATION): // Header location
			{
				var headerLocationData HeaderLocationData
				err = json.Unmarshal(*rawData.Header.Data, &headerLocationData)
				if err != nil {
					return err
				}
				if headerLocationData.Latitude == nil {
					return errors.New("header location latitude is required")
				}
				if headerLocationData.Longitude == nil {
					return errors.New("header location longitude is required")
				}
				if headerLocationData.Name == nil {
					return errors.New("header location name is required")
				}
				if headerLocationData.Address == nil {
					return errors.New("header location address is required")
				}
				t.Header = &TemplateTypeData{
					Type: rawData.Header.Type,
					Data: headerLocationData,
				}
			}
		case string(constants.TEMPLATE_MSG_TYPE_HEADER_AUDIO): // Header audio
			{
				var headerAudioData HeaderAudioData
				err = json.Unmarshal(*rawData.Header.Data, &headerAudioData)
				if err != nil {
					return err
				}
				if headerAudioData.Link == nil {
					return errors.New("header audio link is required")
				}
				t.Header = &TemplateTypeData{
					Type: rawData.Header.Type,
					Data: headerAudioData,
				}
			}
		case string(constants.TEMPLATE_MSG_TYPE_HEADER_VIDEO): // Header video
			{
				var headerVideoData HeaderVideoData
				err = json.Unmarshal(*rawData.Header.Data, &headerVideoData)
				if err != nil {
					return err
				}
				if headerVideoData.Link == nil {
					return errors.New("header video link is required")
				}
				if headerVideoData.Caption == nil {
					return errors.New("header video caption is required")
				}
				t.Header = &TemplateTypeData{
					Type: rawData.Header.Type,
					Data: headerVideoData,
				}
			}
		case string(constants.TEMPLATE_MSG_TYPE_HEADER_DOCUMENT): // Header document
			{
				var headerDocumentData HeaderDocumentData
				err = json.Unmarshal(*rawData.Header.Data, &headerDocumentData)
				if err != nil {
					return err
				}
				if headerDocumentData.Link == nil {
					return errors.New("header document link is required")
				}
				if headerDocumentData.Caption == nil {
					return errors.New("header document caption is required")
				}
				t.Header = &TemplateTypeData{
					Type: rawData.Header.Type,
					Data: headerDocumentData,
				}
			}
		default:
			{
				return errors.New("invalid header type, only image, location, audio, video, document is allowed")
			}
		}
	}

	// Unmarshal body
	// Body is used for text as of now, can add date_time, currency as well ( Not sure how these works so skipped for now )
	if rawData.Body != nil && len(rawData.Body) > 0 {
		for _, rawBody := range rawData.Body {
			if rawBody.Type == nil {
				return errors.New("body type is required")
			}
			if rawBody.Data == nil {
				return errors.New("body data is required")
			}
			switch *rawBody.Type {
			case string(constants.TEMPLATE_MSG_TYPE_TEXT): // Body text
				{
					var bodyTextData BodyTextData
					err = json.Unmarshal(*rawBody.Data, &bodyTextData)
					if err != nil {
						return err
					}
					if bodyTextData.Text == nil {
						return errors.New("body text is required")
					}
					t.Body = append(t.Body, TemplateTypeData{
						Type: rawBody.Type,
						Data: bodyTextData,
					})
				}
			default:
				{
					return errors.New("invalid body type, Only text is allowed")
				}
			}
		}
	}

	// Unmarshal button
	// Button is used for quick reply and url(Otp)
	if rawData.Button != nil && len(rawData.Button) > 0 {
		for _, rawButton := range rawData.Button {
			if rawButton.Type == nil {
				return errors.New("button type is required")
			}
			if rawButton.Data == nil {
				return errors.New("button data is required")
			}
			switch *rawButton.Type {
			case string(constants.TEMPLATE_MSG_TYPE_BUTTON_SUB_TYPE_QUICK_REPLY): // Button quick reply
				{
					var buttonQuickReplyData ButtonQuickReplyData
					err = json.Unmarshal(*rawButton.Data, &buttonQuickReplyData)
					if err != nil {
						return err
					}
					if buttonQuickReplyData.Payload == nil {
						return errors.New("button quick reply payload is required")
					}
					t.Button = append(t.Button, TemplateTypeData{
						Type: rawButton.Type,
						Data: buttonQuickReplyData,
					})
				}
			case string(constants.TEMPLATE_MSG_TYPE_BUTTON_SUB_TYPE_URL): // Button url
				{
					var buttonUrlData ButtonUrlData
					err = json.Unmarshal(*rawButton.Data, &buttonUrlData)
					if err != nil {
						return err
					}
					if buttonUrlData.Text == nil {
						return errors.New("button url text is required")
					}
					t.Button = append(t.Button, TemplateTypeData{
						Type: rawButton.Type,
						Data: buttonUrlData,
					})
				}
			default:
				{
					return errors.New("invalid button type, Only quick reply and url are allowed")
				}
			}
		}
	}

	return nil
}

type TemplatePayloadWhatsapp struct {
	MessagingProduct string                 `json:"messaging_product"`
	RecipientType    string                 `json:"recipient_type"`
	To               string                 `json:"to"`
	Type             string                 `json:"type"`
	Template         map[string]interface{} `json:"template"`
}

type TemplatePayloadClient struct {
	To      string             `json:"to" validate:"required" message:"to is required field"`
	Type    string             `json:"type" validate:"required|TemplateTypeValidator" message:"type is required field and should be 'template' only"`
	Purpose string             `json:"purpose" validate:"required|PurposeValidator" message:"purpose is required and should be either transactional or marketing"`
	Body    TemplateDataStruct `json:"body" validate:"required" message:"body is required field"`
}

func (f TemplatePayloadClient) PurposeValidator(val string) bool {
	return (val == string(constants.MSG_PURPOSE_MARKETING) || val == string(constants.MSG_PURPOSE_TRANSACTIONAL))
}

func (f TemplatePayloadClient) TemplateTypeValidator(val string) bool {
	return val == string(constants.MSG_TYPE_TEMPLATE)
}
