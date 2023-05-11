package types

import (
	"encoding/json"
	"multiBot/constants"
	"multiBot/types/audioUrl"
	"multiBot/types/documentUrl"
	"multiBot/types/imageUrl"
	"multiBot/types/list"
	"multiBot/types/location"
	"multiBot/types/replyButton"
	"multiBot/types/stickerUrl"
	"multiBot/types/template"
	"multiBot/types/text"
	"multiBot/types/videoUrl"
)

type MessagePayloadClient struct {
	To      string      `json:"to" validate:"required" message:"to is required field"`
	Type    string      `json:"type" validate:"required|MessageTypeValidator" message:"type is not valid"`
	Purpose string      `json:"purpose" validate:"required|PurposeValidator" message:"purpose is required and should be either transactional or marketing"`
	Body    interface{} `json:"body" validate:"required" message:"body is required field"`
}

func (f MessagePayloadClient) PurposeValidator(val string) bool {
	return (val == string(constants.MSG_PURPOSE_MARKETING) || val == string(constants.MSG_PURPOSE_TRANSACTIONAL))
}

func (f MessagePayloadClient) MessageTypeValidator(val string) bool {
	switch f.Type {
	case string(constants.MSG_TYPE_TEXT):
		return val == string(constants.MSG_TYPE_TEXT)
	case string(constants.MSG_TYPE_TEMPLATE):
		return val == string(constants.MSG_TYPE_TEMPLATE)
	case string(constants.MSG_TYPE_IMAGE_URL):
		return val == string(constants.MSG_TYPE_IMAGE_URL)
	case string(constants.MSG_TYPE_AUDIO_URL):
		return val == string(constants.MSG_TYPE_AUDIO_URL)
	case string(constants.MSG_TYPE_DOCUMENT_URL):
		return val == string(constants.MSG_TYPE_DOCUMENT_URL)
	case string(constants.MSG_TYPE_VIDEO_URL):
		return val == string(constants.MSG_TYPE_VIDEO_URL)
	case string(constants.MSG_TYPE_STICKER_URL):
		return val == string(constants.MSG_TYPE_STICKER_URL)
	case string(constants.MSG_TYPE_LOCATION):
		return val == string(constants.MSG_TYPE_LOCATION)
	case string(constants.MSG_TYPE_LIST_INTERNAL):
		return val == string(constants.MSG_TYPE_LIST_INTERNAL)
	case string(constants.MSG_TYPE_REPLY_BUTTON_INTERNAL):
		return val == string(constants.MSG_TYPE_REPLY_BUTTON_INTERNAL)
	default:
		return false
	}
}

func (b *MessagePayloadClient) UnmarshalJSON(data []byte) error {
	var rawData struct {
		To      string          `json:"to"`
		Type    string          `json:"type"`
		Purpose string          `json:"purpose"`
		Body    json.RawMessage `json:"body"`
	}
	err := json.Unmarshal(data, &rawData)
	if err != nil {
		return err
	}
	b.To = rawData.To
	b.Type = rawData.Type
	b.Purpose = rawData.Purpose
	switch rawData.Type {
	case string(constants.MSG_TYPE_TEXT):
		{
			var data text.TextDataStruct
			err = json.Unmarshal(rawData.Body, &data)
			if err != nil {
				return err
			}
			b.Body = data
		}
	case string(constants.MSG_TYPE_TEMPLATE):
		{
			var data template.TemplateDataStruct
			err = json.Unmarshal(rawData.Body, &data)
			if err != nil {
				return err
			}
			b.Body = data
		}
	case string(constants.MSG_TYPE_IMAGE_URL):
		{
			var data imageUrl.ImageUrlDataStruct
			err = json.Unmarshal(rawData.Body, &data)
			if err != nil {
				return err
			}
			b.Body = data
		}
	case string(constants.MSG_TYPE_AUDIO_URL):
		{
			var data audioUrl.AudioUrlDataStruct
			err = json.Unmarshal(rawData.Body, &data)
			if err != nil {
				return err
			}
			b.Body = data
		}
	case string(constants.MSG_TYPE_DOCUMENT_URL):
		{
			var data documentUrl.DocumentUrlDataStruct
			err = json.Unmarshal(rawData.Body, &data)
			if err != nil {
				return err
			}
			b.Body = data
		}
	case string(constants.MSG_TYPE_VIDEO_URL):
		{
			var data videoUrl.VideoUrlDataStruct
			err = json.Unmarshal(rawData.Body, &data)
			if err != nil {
				return err
			}
			b.Body = data
		}
	case string(constants.MSG_TYPE_STICKER_URL):
		{
			var data stickerUrl.StickerUrlDataStruct
			err = json.Unmarshal(rawData.Body, &data)
			if err != nil {
				return err
			}
			b.Body = data
		}
	case string(constants.MSG_TYPE_LOCATION):
		{
			var data location.LocationDataStruct
			err = json.Unmarshal(rawData.Body, &data)
			if err != nil {
				return err
			}
			b.Body = data
		}
	case string(constants.MSG_TYPE_LIST_INTERNAL):
		{
			var data list.ListDataStruct
			err = json.Unmarshal(rawData.Body, &data)
			if err != nil {
				return err
			}
			b.Body = data
		}
	case string(constants.MSG_TYPE_REPLY_BUTTON_INTERNAL):
		{
			var data replyButton.ReplyButtonDataStruct
			err = json.Unmarshal(rawData.Body, &data)
			if err != nil {
				return err
			}
			b.Body = data
		}
	}
	return nil
}
