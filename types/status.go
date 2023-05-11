package types

import "multiBot/models"

type StatusPayload struct {
	Ids []string `json:"ids" validate:"required" message:"ids are required"`
}

type StatusMessage struct {
	Id                string                 `json:"id"`
	Body              map[string]interface{} `json:"body"`
	Status            string                 `json:"status"`
	MsgId             string                 `json:"msgId"`
	StatusHistory     []models.StatusHistory `json:"statusHistory"`
	RetryCount        int64                  `json:"retryCount"`
	MaxRetryAvailable int64                  `json:"maxRetryAvailable"`
	Client            string                 `json:"client"`
	Organization      string                 `json:"organization"`
	WhatsappAccount   string                 `json:"whatsappAccount"`
	NetworkType       string                 `json:"msgNetworkType"`
	Purpose           string                 `json:"msgPurpose"`
	CreatedAt         int64                  `json:"createdAt"`
	UpdatedAt         int64                  `json:"updatedAt"`
}

type StatusResponse struct {
	Count    int             `json:"count"`
	Messages []StatusMessage `json:"messages"`
}
