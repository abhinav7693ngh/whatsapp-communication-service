package constants

type StatusType string

const (
	MSG_STATUS_CREATED                               StatusType = "CREATED"
	MSG_STATUS_PROCESSING                            StatusType = "PROCESSING"
	MSG_STATUS_WA_PUSH_SUCCESS                       StatusType = "WA_PUSH_SUCCESS"
	MSG_STATUS_WA_PUSH_FAILED                        StatusType = "WA_PUSH_FAILED"
	MSG_STATUS_WA_PUSH_TIMEOUT                       StatusType = "WA_PUSH_TIMEOUT"
	MSG_STATUS_WA_PUSH_FAILED_INCORRECT_REQUEST      StatusType = "WA_PUSH_FAILED_INCORRECT_REQUEST"
	MSG_STATUS_WA_STATUS_RECEIVED_QUEUE_PUSH_SUCCESS StatusType = "WA_STATUS_RECEIVED_QUEUE_PUSH_SUCCESS"
	MSG_STATUS_WA_STATUS_RECEIVED_QUEUE_PUSH_FAILED  StatusType = "WA_STATUS_RECEIVED_QUEUE_PUSH_FAILED"
	MSG_STATUS_WA_SENT                               StatusType = "WA_SENT"
	MSG_STATUS_WA_FAILED                             StatusType = "WA_FAILED"
	MSG_STATUS_WA_DELIVERED                          StatusType = "WA_DELIVERED"
	MSG_STATUS_WA_UNDELIVERED                        StatusType = "WA_UNDELIVERED"
	MSG_STATUS_WA_READ                               StatusType = "WA_READ"
)

var StatusMapWhatsapp map[string]StatusType = map[string]StatusType{
	"CREATED":                               MSG_STATUS_CREATED,
	"PROCESSING":                            MSG_STATUS_PROCESSING,
	"WA_PUSH_SUCCESS":                       MSG_STATUS_WA_PUSH_SUCCESS,
	"WA_PUSH_FAILED":                        MSG_STATUS_WA_PUSH_FAILED,
	"WA_PUSH_TIMEOUT":                       MSG_STATUS_WA_PUSH_TIMEOUT,
	"WA_PUSH_FAILED_INCORRECT_REQUEST":      MSG_STATUS_WA_PUSH_FAILED_INCORRECT_REQUEST,
	"WA_STATUS_RECEIVED_QUEUE_PUSH_SUCCESS": MSG_STATUS_WA_STATUS_RECEIVED_QUEUE_PUSH_SUCCESS,
	"WA_STATUS_RECEIVED_QUEUE_PUSH_FAILED":  MSG_STATUS_WA_STATUS_RECEIVED_QUEUE_PUSH_FAILED,
	"WA_SENT":                               MSG_STATUS_WA_SENT,
	"WA_FAILED":                             MSG_STATUS_WA_FAILED,
	"WA_DELIVERED":                          MSG_STATUS_WA_DELIVERED,
	"WA_UNDELIVERED":                        MSG_STATUS_WA_UNDELIVERED,
	"WA_READ":                               MSG_STATUS_WA_READ,
}

func CheckValidStatus(status string) bool {
	for k := range StatusMapWhatsapp {
		if k == status {
			return true
		}
	}
	return false
}

var WebhookStatusMap map[string]StatusType = map[string]StatusType{
	"sent":        MSG_STATUS_WA_SENT,
	"failed":      MSG_STATUS_WA_FAILED,
	"delivered":   MSG_STATUS_WA_DELIVERED,
	"undelivered": MSG_STATUS_WA_UNDELIVERED,
	"read":        MSG_STATUS_WA_READ,
}

func CheckWhatsappStatus(status string) bool {
	for k := range WebhookStatusMap {
		if k == status {
			return true
		}
	}
	return false
}
