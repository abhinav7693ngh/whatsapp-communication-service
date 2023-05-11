package template

import (
	"multiBot/constants"
	"multiBot/types/template"
	"strconv"
)

type TemplateLanguage struct {
	Code string `json:"code"`
}

type TemplateParsedBody struct {
	Name       string           `json:"name"`
	Language   TemplateLanguage `json:"language"`
	Components []interface{}    `json:"components"`
}

func makeHeaders(tempHeader template.TemplateTypeData) *CompositeHeader {
	if *tempHeader.Type == "" {
		return nil
	}
	compHeader := CompositeHeader{
		Type: string(constants.TEMPLATE_MSG_TYPE_HEADER),
	}
	switch *tempHeader.Type {
	case string(constants.TEMPLATE_MSG_TYPE_HEADER_IMAGE): // Header image
		{
			compHeader.Parameters = append(compHeader.Parameters, NuclearImage{
				Type: string(constants.TEMPLATE_MSG_TYPE_HEADER_IMAGE),
				Image: ImageObject{
					Link: *tempHeader.Data.(template.HeaderImageData).Link,
				},
			})
		}
	case string(constants.TEMPLATE_MSG_TYPE_HEADER_LOCATION): // Header location
		{
			compHeader.Parameters = append(compHeader.Parameters, NuclearLocation{
				Type: string(constants.TEMPLATE_MSG_TYPE_HEADER_LOCATION),
				Location: LocationObject{
					Latitude:  *tempHeader.Data.(template.HeaderLocationData).Latitude,
					Longitude: *tempHeader.Data.(template.HeaderLocationData).Longitude,
					Name:      *tempHeader.Data.(template.HeaderLocationData).Name,
					Address:   *tempHeader.Data.(template.HeaderLocationData).Address,
				},
			})
		}
	case string(constants.TEMPLATE_MSG_TYPE_HEADER_AUDIO): // Header audio
		{
			compHeader.Parameters = append(compHeader.Parameters, NuclearAudio{
				Type: string(constants.TEMPLATE_MSG_TYPE_HEADER_AUDIO),
				Audio: AudioObject{
					Link: *tempHeader.Data.(template.HeaderAudioData).Link,
				},
			})
		}
	case string(constants.TEMPLATE_MSG_TYPE_HEADER_VIDEO): // Header video
		{
			compHeader.Parameters = append(compHeader.Parameters, NuclearVideo{
				Type: string(constants.TEMPLATE_MSG_TYPE_HEADER_VIDEO),
				Video: VideoObject{
					Link:    *tempHeader.Data.(template.HeaderVideoData).Link,
					Caption: *tempHeader.Data.(template.HeaderVideoData).Caption,
				},
			})
		}
	case string(constants.TEMPLATE_MSG_TYPE_HEADER_DOCUMENT): // Header document
		{
			compHeader.Parameters = append(compHeader.Parameters, NuclearDocument{
				Type: string(constants.TEMPLATE_MSG_TYPE_HEADER_DOCUMENT),
				Document: DocumentObject{
					Link:    *tempHeader.Data.(template.HeaderDocumentData).Link,
					Caption: *tempHeader.Data.(template.HeaderDocumentData).Caption,
				},
			})
		}
	}
	if len(compHeader.Parameters) == 0 {
		return nil
	}
	return &compHeader
}

func makeBody(tempBody []template.TemplateTypeData) *CompositeBody {
	if len(tempBody) == 0 {
		return nil
	}
	compBody := CompositeBody{
		Type: string(constants.TEMPLATE_MSG_TYPE_BODY),
	}
	for _, body := range tempBody { // text body
		switch *body.Type {
		case string(constants.TEMPLATE_MSG_TYPE_TEXT): // text body
			{
				compBody.Parameters = append(compBody.Parameters, NuclearText{
					Type: string(constants.TEMPLATE_MSG_TYPE_TEXT),
					Text: *body.Data.(template.BodyTextData).Text,
				})
			}
		}
	}
	return &compBody
}

func makeButtons(tempButton []template.TemplateTypeData) []CompositeButton {
	compButtons := []CompositeButton{}
	for index, button := range tempButton {
		compButton := CompositeButton{
			Type:  string(constants.TEMPLATE_MSG_TYPE_BUTTON),
			Index: strconv.Itoa(index),
		}
		switch *button.Type {
		case string(constants.TEMPLATE_MSG_TYPE_BUTTON_SUB_TYPE_QUICK_REPLY): // reply button
			{
				compButton.Sub_Type = string(constants.TEMPLATE_MSG_TYPE_BUTTON_SUB_TYPE_QUICK_REPLY)
				compButton.Parameters = append(compButton.Parameters, NuclearPayload{
					Type:    string(constants.TEMPLATE_MSG_TYPE_BUTTON_PAYLOAD),
					Payload: *button.Data.(template.ButtonQuickReplyData).Payload,
				})
			}
		case string(constants.TEMPLATE_MSG_TYPE_BUTTON_SUB_TYPE_URL): // otp button
			{
				compButton.Sub_Type = string(constants.TEMPLATE_MSG_TYPE_BUTTON_SUB_TYPE_URL)
				compButton.Parameters = append(compButton.Parameters, NuclearOtpText{
					Type: string(constants.TEMPLATE_MSG_TYPE_TEXT),
					Text: *button.Data.(template.ButtonUrlData).Text,
				})
			}
		}
		compButtons = append(compButtons, compButton)
	}
	return compButtons
}

func ParseTemplate(temp template.TemplateDataStruct) TemplateParsedBody {
	var madeHeaders *CompositeHeader
	if temp.Header != nil {
		madeHeaders = makeHeaders(*temp.Header)
	}
	madeBody := makeBody(temp.Body)
	madeButtons := makeButtons(temp.Button)
	parsedTemp := TemplateParsedBody{
		Name:       temp.Name,
		Language:   TemplateLanguage{Code: temp.Language},
		Components: nil,
	}
	if madeHeaders != nil {
		parsedTemp.Components = append(parsedTemp.Components, *madeHeaders)
	}
	if madeBody != nil {
		parsedTemp.Components = append(parsedTemp.Components, *madeBody)
	}
	for _, button := range madeButtons {
		parsedTemp.Components = append(parsedTemp.Components, button)
	}
	return parsedTemp
}
