package replyButton

// ================ Reply ================ //

type ReplyObject struct {
	Title string `json:"title"`
	Id    string `json:"id"`
}

type NuclearReply struct {
	Type  string      `json:"type"`
	Reply ReplyObject `json:"location"`
}
