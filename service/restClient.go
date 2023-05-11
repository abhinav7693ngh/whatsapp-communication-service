package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"multiBot/config"
	"multiBot/errorCodes"
	"multiBot/models"
	"net/http"
	"time"
)

type RestClientWhatsapp struct {
	phoneId          string
	accessToken      string
	apiTimeout       time.Duration
	retryCount       int
	retryWaitTime    time.Duration
	retryMaxWaitTime time.Duration
	pathParams       map[string]string
}

func InitRestClient(params map[string]interface{}, message models.Message) (rc *RestClientWhatsapp) {

	cfg := config.GetConfig()

	waAccountInfo := config.GetWhatsappAccountInfo(message.WhatsappAccountIdentifier)

	pathParam := make(map[string]string)
	for key, val := range params {
		pathParam[key] = val.(string)
	}

	rc = &RestClientWhatsapp{
		phoneId:          waAccountInfo.PHONE_ID,
		accessToken:      waAccountInfo.ACCESS_TOKEN,
		apiTimeout:       cfg.WHATSAPP_API.TIMEOUT_SECONDS,
		retryCount:       cfg.WHATSAPP_API.RETRY_COUNT,
		retryWaitTime:    cfg.WHATSAPP_API.RETRY_WAIT_TIME * time.Second,
		retryMaxWaitTime: cfg.WHATSAPP_API.RETRY_MAX_WAIT_TIME * time.Second,
		pathParams:       pathParam,
	}

	return
}

// TODO: Add retry and debug mode
func (rcw *RestClientWhatsapp) Post(body *[]byte) (*http.Response, error) {
	cfg := config.GetConfig()

	postUrl := fmt.Sprintf("%s/%s/%s",
		cfg.WHATSAPP_API.BASE_URL,
		rcw.phoneId,
		cfg.WHATSAPP_API.POST_ENDPOINT,
	)

	ctx, cancel := context.WithTimeout(context.Background(), rcw.apiTimeout*time.Second)
	defer cancel()

	sendReq, err := http.NewRequestWithContext(ctx, "POST", postUrl, bytes.NewBuffer(*body))
	if err != nil {
		return nil, errors.New(errorCodes.ERROR_CREATING_NEW_REQUEST)
	}
	sendReq.Header.Set("Authorization", "Bearer "+rcw.accessToken)
	sendReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	sendResp, err := client.Do(sendReq)
	if err != nil {
		return nil, err
	}

	return sendResp, nil
}
