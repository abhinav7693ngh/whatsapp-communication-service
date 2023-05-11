package constants

/*

===> type mapping

image --> imageUrl, imageId
document --> documentUrl, documentId
sticker --> stickerUrl, stickerId
video --> videoUrl, videoId
audio --> audioUrl, audioId
interactive --> list, button, product, product_list

*/

type MsgType string
type InteractiveType string
type ImageType string
type AudioType string
type DocumentType string
type StickerType string
type VideoType string

type MsgPurposeType string
type MsgNetworkType string

const (
	MSG_TYPE_TEXT        MsgType = "text"
	MSG_TYPE_IMAGE       MsgType = "image"
	MSG_TYPE_TEMPLATE    MsgType = "template"
	MSG_TYPE_AUDIO       MsgType = "audio"
	MSG_TYPE_DOCUMENT    MsgType = "document"
	MSG_TYPE_STICKER     MsgType = "sticker"
	MSG_TYPE_VIDEO       MsgType = "video"
	MSG_TYPE_LOCATION    MsgType = "location"
	MSG_TYPE_INTERACTIVE MsgType = "interactive"
)

const (
	MSG_TYPE_IMAGE_URL ImageType = "imageUrl"
)

const (
	MSG_TYPE_DOCUMENT_URL ImageType = "documentUrl"
)

const (
	MSG_TYPE_STICKER_URL ImageType = "stickerUrl"
)

const (
	MSG_TYPE_VIDEO_URL ImageType = "videoUrl"
)

const (
	MSG_TYPE_AUDIO_URL ImageType = "audioUrl"
)

const (
	MSG_TYPE_LIST_INTERNAL            InteractiveType = "list"
	MSG_TYPE_LIST_INTERACTIVE         InteractiveType = "list"
	MSG_TYPE_REPLY_BUTTON_INTERNAL    InteractiveType = "replyButton"
	MSG_TYPE_REPLY_BUTTON_INTERACTIVE InteractiveType = "button"
)

const (
	MSG_TYPE_WA_WEBHOOK_TEXT MsgType = "text"
)

// ============================================================================ //

const (
	MSG_PURPOSE_MARKETING     MsgPurposeType = "marketing"
	MSG_PURPOSE_TRANSACTIONAL MsgPurposeType = "transactional"
	MSG_PURPOSE_CHAT          MsgPurposeType = "chat"
)

const (
	MSG_NETWORK_OUTGOING MsgNetworkType = "OUTGOING"
	MSG_NETWORK_INCOMING MsgNetworkType = "INCOMING"
)
