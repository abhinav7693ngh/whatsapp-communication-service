package replyButton

import (
	"multiBot/constants"
	"multiBot/types/replyButton"

	"github.com/google/uuid"
)

type ReplyButtonBody struct {
	Text string `json:"text"`
}

type ReplyButtonAction struct {
	Buttons []NuclearReply `json:"buttons"`
}

type ReplyButtonParsedBody struct {
	Type   string            `json:"type"`
	Body   ReplyButtonBody   `json:"body"`
	Action ReplyButtonAction `json:"action"`
}

func (r *ReplyButtonParsedBody) fill_defaults() {
	r.Type = string(constants.MSG_TYPE_REPLY_BUTTON_INTERACTIVE)
}

func makeActionButton(buttons []string) ReplyButtonAction {
	var action ReplyButtonAction
	// all buttons considered as reply buttons, change here if you want to add more types
	for _, button := range buttons {
		action.Buttons = append(action.Buttons, NuclearReply{
			Type: string(constants.REPLY_BUTTON_MSG_TYPE_BUTTON_REPLY),
			Reply: ReplyObject{
				Title: button,
				Id:    uuid.New().String(),
			},
		})
	}
	return action
}

func ParseReplyButton(replies replyButton.ReplyButtonDataStruct) ReplyButtonParsedBody {
	var replyButtonParsedBody ReplyButtonParsedBody
	replyButtonParsedBody.fill_defaults()
	replyButtonParsedBody.Body.Text = replies.Body
	replyButtonParsedBody.Action = makeActionButton(replies.Buttons)
	return replyButtonParsedBody
}
