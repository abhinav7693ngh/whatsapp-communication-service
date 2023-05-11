package constants

type TemplateMsgType string

const (
	TEMPLATE_MSG_TYPE_TEXT                        TemplateMsgType = "text"
	TEMPLATE_MSG_TYPE_HEADER                      TemplateMsgType = "header"
	TEMPLATE_MSG_TYPE_HEADER_LOCATION             TemplateMsgType = "location"
	TEMPLATE_MSG_TYPE_HEADER_IMAGE                TemplateMsgType = "image"
	TEMPLATE_MSG_TYPE_HEADER_DOCUMENT             TemplateMsgType = "document"
	TEMPLATE_MSG_TYPE_HEADER_AUDIO                TemplateMsgType = "audio"
	TEMPLATE_MSG_TYPE_HEADER_VIDEO                TemplateMsgType = "video"
	TEMPLATE_MSG_TYPE_BUTTON                      TemplateMsgType = "button"
	TEMPLATE_MSG_TYPE_BUTTON_PAYLOAD              TemplateMsgType = "payload"
	TEMPLATE_MSG_TYPE_BUTTON_SUB_TYPE_QUICK_REPLY TemplateMsgType = "quick_reply"
	TEMPLATE_MSG_TYPE_BUTTON_SUB_TYPE_URL         TemplateMsgType = "url"
	TEMPLATE_MSG_TYPE_BODY                        TemplateMsgType = "body"

	// NOT IN USE CURRENTLY AS PER USAGE
	TEMPLATE_MSG_TYPE_CURRENCY  TemplateMsgType = "currency"
	TEMPLATE_MSG_TYPE_DATE_TIME TemplateMsgType = "date_time"
)
