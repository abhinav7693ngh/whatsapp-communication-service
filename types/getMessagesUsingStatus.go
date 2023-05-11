package types

import "multiBot/models"

type GetMessagesUsingStatusPayload struct {
	Status string `json:"status" validate:"required" message:"status is required"`
}

type MessageUsingStatus struct {
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

type GetMessagesUsingStatusResponse struct {
	Page     int                  `json:"page"`
	Limit    int                  `json:"limit"`
	Total    int64                `json:"total"`
	Count    int                  `json:"count"`
	Messages []MessageUsingStatus `json:"messages"`
}
