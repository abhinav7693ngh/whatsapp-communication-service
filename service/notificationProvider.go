package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"multiBot/constants"
	"multiBot/errorCodes"
	"multiBot/logger"
	"multiBot/models"
	"net/http"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WhatsAppResponseMessage struct {
	Id string `json:"id"`
}

type WhatsAppResponse struct {
	Messages []WhatsAppResponseMessage `json:"messages"`
}

type ErrorStruct struct {
	Message string `json:"message"`
}

type WhatsappError struct {
	Error ErrorStruct `json:"error"`
}

type NotificationProviderWhatsapp struct {
}

func (npw *NotificationProviderWhatsapp) prepareMessage(message models.Message) (*[]byte, error) {
	switch string(message.InternalMsgType) {
	case string(constants.MSG_TYPE_TEXT):
		whatsappJSON, err := PrepareTextTypeMessage(message)
		if err != nil {
			return nil, err
		}
		return whatsappJSON, nil
	case string(constants.MSG_TYPE_TEMPLATE):
		whatsappJSON, err := PrepareTemplateTypeMessage(message)
		if err != nil {
			return nil, err
		}
		return whatsappJSON, nil
	case string(constants.MSG_TYPE_IMAGE_URL):
		whatsappJSON, err := PrepareImageUrlTypeMessage(message)
		if err != nil {
			return nil, err
		}
		return whatsappJSON, nil
	case string(constants.MSG_TYPE_AUDIO_URL):
		whatsappJSON, err := PrepareAudioUrlTypeMessage(message)
		if err != nil {
			return nil, err
		}
		return whatsappJSON, nil
	case string(constants.MSG_TYPE_DOCUMENT_URL):
		whatsappJSON, err := PrepareDocumentUrlTypeMessage(message)
		if err != nil {
			return nil, err
		}
		return whatsappJSON, nil
	case string(constants.MSG_TYPE_STICKER_URL):
		whatsappJSON, err := PrepareStickerUrlTypeMessage(message)
		if err != nil {
			return nil, err
		}
		return whatsappJSON, nil
	case string(constants.MSG_TYPE_VIDEO_URL):
		whatsappJSON, err := PrepareVideoUrlTypeMessage(message)
		if err != nil {
			return nil, err
		}
		return whatsappJSON, nil
	case string(constants.MSG_TYPE_LOCATION):
		whatsappJSON, err := PrepareLocationTypeMessage(message)
		if err != nil {
			return nil, err
		}
		return whatsappJSON, nil
	case string(constants.MSG_TYPE_LIST_INTERNAL):
		whatsappJSON, err := PrepareListTypeMessage(message)
		if err != nil {
			return nil, err
		}
		return whatsappJSON, nil
	case string(constants.MSG_TYPE_REPLY_BUTTON_INTERNAL):
		whatsappJSON, err := PrepareReplyButtonTypeMessage(message)
		if err != nil {
			return nil, err
		}
		return whatsappJSON, nil
	default:
		{
			return nil, errors.New("SendWhatsappMessage GoRoutine: Default case in prepare message, id: " + message.Id.Hex() + ", error: message type is not valid")
		}
	}
}

func (npw *NotificationProviderWhatsapp) ProcessWhatsappError(errType string, err error, filter primitive.M, message models.Message) {
	var msgs [2]string
	switch string(errType) {
	case "prepare_message_err":
		logger.LogError(nil, err.Error(), nil)

		msgs[0] = "SendWhatsappMessage GoRoutine: Error in updating message status ( prepare message error ) with id: "
		msgs[1] = "SendWhatsappMessage GoRoutine: Updated message status ( prepare message error ) with id: "

		StatusUpdateWhatsappPushFailed(msgs, filter, message)
		return
	case "nil_send_resp_err":
		logger.LogError(
			nil,
			"SendWhatsappMessage GoRoutine: error creating new request, id: "+message.Id.Hex()+", error: "+err.Error(),
			nil,
		)

		msgs[0] = "SendWhatsappMessage GoRoutine: Error in updating message status ( error creating new request ) with id: "
		msgs[1] = "SendWhatsappMessage GoRoutine: Updated message status ( error creating new request ) with id: "

		StatusUpdateWhatsappPushFailed(msgs, filter, message)
		return
	case "send_resp_err":
		if e, ok := err.(*url.Error); ok {
			if e.Err == context.DeadlineExceeded {
				possible := message.PossibleToUpdateStatus(constants.MSG_STATUS_WA_PUSH_TIMEOUT, "SendWhatsappMessage")
				if !possible {
					return
				}
				currentTime := time.Now().UnixMilli()
				update := bson.M{
					"$set": bson.M{
						"status":    constants.MSG_STATUS_WA_PUSH_TIMEOUT,
						"updatedAt": currentTime,
					},
					"$push": bson.M{
						"statusHistory": bson.M{
							"status":    constants.MSG_STATUS_WA_PUSH_TIMEOUT,
							"updatedAt": currentTime,
						},
					},
					"$inc": bson.M{
						"retryCount": 1,
					},
				}

				updatedNo, err := message.UpdateOneMessage(filter, update)
				if err != nil {
					logger.LogError(
						nil,
						"SendWhatsappMessage GoRoutine: Error in updating message in deadline case with id: "+message.Id.Hex()+", error: "+err.Error(),
						nil,
					)
				} else if updatedNo != nil {
					if *updatedNo > 0 {
						logger.LogInfo(
							nil,
							"SendWhatsappMessage GoRoutine: Updated message in deadline case with id: "+message.Id.Hex(),
							nil,
						)
					}
				}
				return
			}

			msgs[0] = "SendWhatsappMessage GoRoutine: Error in updating message status ( url error deadline ) with id: "
			msgs[1] = "SendWhatsappMessage GoRoutine: Updated message status ( url error deadline ) with id: "

			StatusUpdateWhatsappPushFailed(msgs, filter, message)
			return
		}

		msgs[0] = "SendWhatsappMessage GoRoutine: Error in updating message status ( url error deadline ) with id: "
		msgs[1] = "SendWhatsappMessage GoRoutine: Updated message status ( url error deadline ) with id: "

		StatusUpdateWhatsappPushFailed(msgs, filter, message)
		return
	case "status_ok_err":
		logger.LogError(
			nil,
			"SendWhatsappMessage GoRoutine: error in reading status OK response, id: "+message.Id.Hex()+", error: "+err.Error(),
			nil,
		)

		msgs[0] = "SendWhatsappMessage GoRoutine: Error in updating message in status OK response ( cannot read response body ) with id: "
		msgs[1] = "SendWhatsappMessage GoRoutine: Updated message in status OK response ( cannot read response body ) with id: "

		StatusUpdateWhatsappPushFailed(msgs, filter, message)
		return
	case "unmarshal_WAbody_err":
		logger.LogError(
			nil,
			"SendWhatsappMessage GoRoutine: error in json unmarshaling in status OK case, id: "+message.Id.Hex()+", error: "+err.Error(),
			nil,
		)

		msgs[0] = "SendWhatsappMessage GoRoutine: Error in updating message in status OK response ( json unmarshal error ) with id: "
		msgs[1] = "SendWhatsappMessage GoRoutine: Updated message in status OK response ( json unmarshal error ) with id: "

		StatusUpdateWhatsappPushFailed(msgs, filter, message)
		return
	case "wa_messages_length_err":
		logger.LogInfo(
			nil,
			"SendWhatsappMessage GoRoutine: Got length zero in whatsapp messages, id: "+message.Id.Hex(),
			nil,
		)

		msgs[0] = "SendWhatsappMessage GoRoutine: Error in updating message in status OK response ( reponse length 0 ) with id: "
		msgs[1] = "SendWhatsappMessage GoRoutine: Updated message in status OK response ( response length 0 ) with id: "

		StatusUpdateWhatsappPushFailed(msgs, filter, message)
		return
	case "status_bad_req_err":
		logger.LogError(
			nil,
			"SendWhatsappMessage GoRoutine: error in reading bad request response, id: "+message.Id.Hex()+", error: "+err.Error(),
			nil,
		)

		msgs[0] = "SendWhatsappMessage GoRoutine: Error in updating message in bad request response ( cannot read response body ) with id: "
		msgs[1] = "SendWhatsappMessage GoRoutine: Updated message in bad request response ( cannot read response body ) with id: "

		StatusUpdateWhatsappPushFailed(msgs, filter, message)
		return
	case "unmarshal_WAerr_err":
		logger.LogError(
			nil,
			"SendWhatsappMessage GoRoutine: error in json unmarshaling in bad request case, id: "+message.Id.Hex()+", error: "+err.Error(),
			nil,
		)

		msgs[0] = "SendWhatsappMessage GoRoutine: Error in updating message in bad request response ( json unmarshal error ) with id: "
		msgs[1] = "SendWhatsappMessage GoRoutine: Updated message in bad request response ( json unmarshal error ) with id: "

		StatusUpdateWhatsappPushFailed(msgs, filter, message)
		return
	}
}

func (npw *NotificationProviderWhatsapp) SendWhatsappMessage(message models.Message) {
	var msgs [2]string

	filter := bson.M{
		"_id": message.Id,
	}

	whatsappJSON, err := npw.prepareMessage(message)
	if err != nil {
		npw.ProcessWhatsappError("prepare_message_err", err, filter, message)
		return
	}

	var params map[string]interface{}
	client := InitRestClient(params, message)

	sendResp, err := client.Post(whatsappJSON)
	if err != nil {
		if err.Error() == errorCodes.ERROR_CREATING_NEW_REQUEST {
			npw.ProcessWhatsappError("nil_send_resp_err", err, filter, message)
			return
		}
		npw.ProcessWhatsappError("send_resp_err", err, filter, message)
		return
	}
	defer sendResp.Body.Close()

	var whatsappResp WhatsAppResponse
	if sendResp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(sendResp.Body)
		if err != nil {
			npw.ProcessWhatsappError("status_ok_err", err, filter, message)
			return
		}

		err = json.Unmarshal(bodyBytes, &whatsappResp)
		if err != nil {
			npw.ProcessWhatsappError("unmarshal_WAbody_err", err, filter, message)
			return
		}

		if len(whatsappResp.Messages) > 0 {
			possible := message.PossibleToUpdateStatus(constants.MSG_STATUS_WA_PUSH_SUCCESS, "SendWhatsappMessage GoRoutine")
			if !possible {
				return
			}
			currentTime := time.Now().UnixMilli()
			update := bson.M{
				"$set": bson.M{
					"msgId":     whatsappResp.Messages[0].Id,
					"status":    constants.MSG_STATUS_WA_PUSH_SUCCESS,
					"updatedAt": currentTime,
				},
				"$push": bson.M{
					"statusHistory": bson.M{
						"status":    constants.MSG_STATUS_WA_PUSH_SUCCESS,
						"updatedAt": currentTime,
					},
				},
				"$inc": bson.M{
					"retryCount": 1,
				},
			}

			updatedNo, err := message.UpdateOneMessage(filter, update)
			if err != nil {
				logger.LogError(
					nil,
					"SendWhatsappMessage GoRoutine: Error in updating message in status OK case with id: "+message.Id.Hex()+", error: "+err.Error(),
					nil,
				)
			} else if updatedNo != nil {
				if *updatedNo > 0 {
					logger.LogInfo(
						nil,
						"SendWhatsappMessage GoRoutine: Updated message in status OK case with id: "+message.Id.Hex(),
						nil,
					)
				}
			}

		} else {
			npw.ProcessWhatsappError("wa_messages_length_err", err, filter, message)
			return
		}

	} else if sendResp.StatusCode == http.StatusBadRequest {
		bodyBytes, err := io.ReadAll(sendResp.Body)
		if err != nil {
			npw.ProcessWhatsappError("status_bad_req_err", err, filter, message)
			return
		}

		var whatsappErr WhatsappError

		err = json.Unmarshal(bodyBytes, &whatsappErr)
		if err != nil {
			npw.ProcessWhatsappError("unmarshal_WAerr_err", err, filter, message)
			return
		} else {
			logger.LogError(
				nil,
				"SendWhatsappMessage GoRoutine: error in whatsapp response in bad request case, id: "+message.Id.Hex()+", error: "+whatsappErr.Error.Message,
				nil,
			)
		}

		possible := message.PossibleToUpdateStatus(constants.MSG_STATUS_WA_PUSH_FAILED_INCORRECT_REQUEST, "SendWhatsappMessage GoRoutine")
		if !possible {
			return
		}

		currentTime := time.Now().UnixMilli()
		update := bson.M{
			"$set": bson.M{
				"status":    constants.MSG_STATUS_WA_PUSH_FAILED_INCORRECT_REQUEST,
				"updatedAt": currentTime,
			},
			"$push": bson.M{
				"statusHistory": bson.M{
					"status":    constants.MSG_STATUS_WA_PUSH_FAILED_INCORRECT_REQUEST,
					"updatedAt": currentTime,
				},
			},
			"$inc": bson.M{
				"retryCount": 1,
			},
		}

		updatedNo, err := message.UpdateOneMessage(filter, update)
		if err != nil {
			logger.LogError(
				nil,
				"SendWhatsappMessage GoRoutine: Error in updating message in bad request case with id: "+message.Id.Hex()+", error: "+err.Error(),
				nil,
			)
		} else if updatedNo != nil {
			if *updatedNo > 0 {
				logger.LogInfo(
					nil,
					"SendWhatsappMessage GoRoutine: Updated message status in bad request case with id: "+message.Id.Hex(),
					nil,
				)
			}
		}

	} else {
		bodyBytes, err := io.ReadAll(sendResp.Body)
		if err != nil {
			logger.LogError(
				nil,
				"SendWhatsappMessage GoRoutine: error in reading bad request response, id: "+message.Id.Hex()+", error: "+err.Error(),
				nil,
			)
		}

		var whatsappErr WhatsappError

		err = json.Unmarshal(bodyBytes, &whatsappErr)
		if err != nil {
			logger.LogError(
				nil,
				"SendWhatsappMessage GoRoutine: error in json unmarshaling in last else case, id: "+message.Id.Hex()+", error: "+err.Error(),
				nil,
			)
		} else {
			logger.LogError(
				nil,
				"SendWhatsappMessage GoRoutine: error in whatsapp response in last else case, id: "+message.Id.Hex()+", error: "+whatsappErr.Error.Message,
				nil,
			)
		}

		msgs[0] = "SendWhatsappMessage GoRoutine: Error in updating message in last else case with id: "
		msgs[1] = "SendWhatsappMessage GoRoutine: Updated message in last else case with id: "

		StatusUpdateWhatsappPushFailed(msgs, filter, message)
		return
	}
}
