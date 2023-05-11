package models

import (
	"errors"
	"multiBot/config"
	"multiBot/constants"
	"multiBot/database"
	"multiBot/errorCodes"
	"multiBot/logger"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var stateMachine map[constants.StatusType][]constants.StatusType = map[constants.StatusType][]constants.StatusType{
	constants.MSG_STATUS_CREATED: {
		constants.MSG_STATUS_PROCESSING,
		constants.MSG_STATUS_WA_PUSH_SUCCESS,
		constants.MSG_STATUS_WA_PUSH_FAILED,
		constants.MSG_STATUS_WA_PUSH_TIMEOUT,
		constants.MSG_STATUS_WA_PUSH_FAILED_INCORRECT_REQUEST,
		constants.MSG_STATUS_WA_STATUS_RECEIVED_QUEUE_PUSH_SUCCESS,
		constants.MSG_STATUS_WA_STATUS_RECEIVED_QUEUE_PUSH_FAILED,
		constants.MSG_STATUS_WA_SENT,
		constants.MSG_STATUS_WA_DELIVERED,
		constants.MSG_STATUS_WA_FAILED,
		constants.MSG_STATUS_WA_UNDELIVERED,
		constants.MSG_STATUS_WA_READ,
	},
	constants.MSG_STATUS_PROCESSING: {
		constants.MSG_STATUS_WA_PUSH_SUCCESS,
		constants.MSG_STATUS_WA_PUSH_FAILED,
		constants.MSG_STATUS_WA_PUSH_TIMEOUT,
		constants.MSG_STATUS_WA_PUSH_FAILED_INCORRECT_REQUEST,
		constants.MSG_STATUS_WA_STATUS_RECEIVED_QUEUE_PUSH_SUCCESS,
		constants.MSG_STATUS_WA_STATUS_RECEIVED_QUEUE_PUSH_FAILED,
		constants.MSG_STATUS_WA_SENT,
		constants.MSG_STATUS_WA_DELIVERED,
		constants.MSG_STATUS_WA_FAILED,
		constants.MSG_STATUS_WA_UNDELIVERED,
		constants.MSG_STATUS_WA_READ,
	},
	constants.MSG_STATUS_WA_PUSH_SUCCESS: {
		constants.MSG_STATUS_WA_STATUS_RECEIVED_QUEUE_PUSH_SUCCESS,
		constants.MSG_STATUS_WA_STATUS_RECEIVED_QUEUE_PUSH_FAILED,
		constants.MSG_STATUS_WA_SENT,
		constants.MSG_STATUS_WA_DELIVERED,
		constants.MSG_STATUS_WA_FAILED,
		constants.MSG_STATUS_WA_UNDELIVERED,
		constants.MSG_STATUS_WA_READ,
	},
	constants.MSG_STATUS_WA_PUSH_FAILED: {
		constants.MSG_STATUS_PROCESSING,
		constants.MSG_STATUS_WA_PUSH_SUCCESS,
		constants.MSG_STATUS_WA_PUSH_TIMEOUT,
		constants.MSG_STATUS_WA_PUSH_FAILED_INCORRECT_REQUEST,
		constants.MSG_STATUS_WA_STATUS_RECEIVED_QUEUE_PUSH_SUCCESS,
		constants.MSG_STATUS_WA_STATUS_RECEIVED_QUEUE_PUSH_FAILED,
		constants.MSG_STATUS_WA_SENT,
		constants.MSG_STATUS_WA_DELIVERED,
		constants.MSG_STATUS_WA_FAILED,
		constants.MSG_STATUS_WA_UNDELIVERED,
		constants.MSG_STATUS_WA_READ,
	},
	constants.MSG_STATUS_WA_STATUS_RECEIVED_QUEUE_PUSH_SUCCESS: {
		constants.MSG_STATUS_WA_SENT,
		constants.MSG_STATUS_WA_DELIVERED,
		constants.MSG_STATUS_WA_FAILED,
		constants.MSG_STATUS_WA_UNDELIVERED,
		constants.MSG_STATUS_WA_READ,
	},
	constants.MSG_STATUS_WA_STATUS_RECEIVED_QUEUE_PUSH_FAILED: {
		constants.MSG_STATUS_WA_STATUS_RECEIVED_QUEUE_PUSH_SUCCESS,
	},
	constants.MSG_STATUS_WA_SENT: {
		constants.MSG_STATUS_WA_DELIVERED,
		constants.MSG_STATUS_WA_FAILED,
		constants.MSG_STATUS_WA_UNDELIVERED,
		constants.MSG_STATUS_WA_READ,
	},
	constants.MSG_STATUS_WA_DELIVERED: {
		constants.MSG_STATUS_WA_READ,
	},
	constants.MSG_STATUS_WA_FAILED:                        {},
	constants.MSG_STATUS_WA_UNDELIVERED:                   {},
	constants.MSG_STATUS_WA_PUSH_TIMEOUT:                  {},
	constants.MSG_STATUS_WA_PUSH_FAILED_INCORRECT_REQUEST: {},
	constants.MSG_STATUS_WA_READ:                          {},
}

type SenderReceiverStruct struct {
	Contact string `json:"contact" bson:"contact"`
}

type StatusHistory struct {
	Status    constants.StatusType `json:"status" bson:"status"`
	UpdatedAt int64                `json:"updatedAt" bson:"updatedAt"`
}

type Message struct {
	Id                        primitive.ObjectID     `json:"_id,omitempty" bson:"_id,omitempty"`
	Body                      map[string]interface{} `json:"body" bson:"body"`
	Sender                    SenderReceiverStruct   `json:"sender" bson:"sender"`
	Receiver                  SenderReceiverStruct   `json:"receiver" bson:"receiver"`
	MsgId                     string                 `json:"msgId" bson:"msgId"`
	Status                    constants.StatusType   `json:"status" bson:"status"`
	MsgType                   constants.MsgType      `json:"msgType" bson:"msgType"`
	InternalMsgType           constants.MsgType      `json:"InternalType" bson:"internalType"`
	StatusHistory             []StatusHistory        `json:"statusHistory" bson:"statusHistory"`
	RetryCount                int64                  `json:"retryCount" bson:"retryCount"`
	MaxRetryAvailable         int64                  `json:"maxRetryAvailable" bson:"maxRetryAvailable"`
	ClientIdentifier          string                 `json:"clientIdentifier" bson:"clientIdentifier"`
	WhatsappAccountIdentifier string                 `json:"whatsappAccountIdentifier" bson:"whatsappAccountIdentifier"`
	UserContact               string                 `json:"userContact" bson:"userContact"`
	MsgNetworkType            string                 `json:"msgNetworkType" bson:"msgNetworkType"`
	MsgPurpose                string                 `json:"msgPurpose" bson:"msgPurpose"`
	CreatedAt                 int64                  `json:"createdAt" bson:"createdAt"`
	UpdatedAt                 int64                  `json:"updatedAt" bson:"updatedAt"`
}

func GetNewMessageStructForSending( // for sending purpose thats why outgoing
	body map[string]interface{},
	clientIdentifier string,
	waAccountIdentifier string,
	msgPurpose string,
	receiver string,
	sender string,
	status constants.StatusType,
	msgType constants.MsgType,
	internalMsgType constants.MsgType,
	statusHistory []StatusHistory,
) Message {
	cfg := config.GetConfig()
	return Message{
		Body: body,
		Sender: SenderReceiverStruct{
			Contact: sender,
		},
		Receiver: SenderReceiverStruct{
			Contact: receiver,
		},
		MsgId:                     "",
		Status:                    status,
		MsgType:                   msgType,
		InternalMsgType:           internalMsgType,
		RetryCount:                -1,
		MaxRetryAvailable:         int64(cfg.WHATSAPP_CONSUMER.MAX_WHATSAPP_PUSH_RETRY_AVAILABLE),
		StatusHistory:             statusHistory,
		ClientIdentifier:          clientIdentifier,
		WhatsappAccountIdentifier: waAccountIdentifier,
		UserContact:               receiver,
		MsgNetworkType:            string(constants.MSG_NETWORK_OUTGOING),
		MsgPurpose:                msgPurpose,
		CreatedAt:                 time.Now().UnixMilli(),
		UpdatedAt:                 time.Now().UnixMilli(),
	}
}

func (m Message) CountDocuments(filter primitive.M, options *options.FindOptions) (*int64, error) {
	mg := database.GetDB()
	total, err := mg.MessageCollection.CountDocuments(database.GetDB().MongoCtx, filter)
	if err != nil {
		return nil, err
	}
	if &total != nil {
		return &total, nil
	}
	return nil, errors.New(errorCodes.NO_RECORD_FOUND)
}

func (m Message) PossibleToUpdateStatus(nextStatus constants.StatusType, source string) bool {
	currentStatus := m.Status
	nextPossibleStates := stateMachine[currentStatus]
	for _, val := range nextPossibleStates {
		if val == nextStatus {
			return true
		}
	}
	logger.LogError(
		nil,
		source+": Status update not possible from "+string(m.Status)+" to "+string(nextStatus)+" for msg with id: "+m.Id.Hex(),
		nil,
	)
	return false
}

func (m Message) InsertOneMessage(message Message) (*string, error) {
	mg := database.GetDB()
	inserted, err := mg.MessageCollection.InsertOne(database.GetDB().MongoCtx, message)
	if err != nil {
		return nil, err
	}
	objectIdString := GetStringObjectIdFromBson(inserted.InsertedID)
	if &objectIdString != nil {
		return &objectIdString, nil
	}
	return nil, errors.New(errorCodes.NO_RECORD_FOUND)
}

func (m Message) InsertManyMessages(messages []Message) (*[]string, error) {
	mg := database.GetDB()
	var msgs []interface{}
	for _, msg := range messages {
		msgs = append(msgs, msg)
	}
	insertedAll, err := mg.MessageCollection.InsertMany(database.GetDB().MongoCtx, msgs)
	if err != nil {
		return nil, err
	}
	var objectIds []string
	for _, id := range insertedAll.InsertedIDs {
		objectIds = append(objectIds, GetStringObjectIdFromBson(id))
	}
	if &objectIds != nil {
		return &objectIds, nil
	}
	return nil, errors.New(errorCodes.NO_RECORD_FOUND)
}

func (m Message) UpdateOneMessage(filter primitive.M, update primitive.M) (*int64, error) {
	mg := database.GetDB()
	result, err := mg.MessageCollection.UpdateOne(database.GetDB().MongoCtx, filter, update)
	if err != nil {
		return nil, err
	}
	if &result != nil {
		return &result.ModifiedCount, nil
	}
	return nil, errors.New(errorCodes.NO_RECORD_FOUND)
}

func (m Message) UpdateManyMessages(filter primitive.M, update primitive.M) (*int64, error) {
	mg := database.GetDB()
	result, err := mg.MessageCollection.UpdateMany(database.GetDB().MongoCtx, filter, update)
	if err != nil {
		return nil, err
	}
	if &result != nil {
		return &result.ModifiedCount, nil
	}
	return nil, errors.New(errorCodes.NO_RECORD_FOUND)
}

func (m Message) DeleteOneMessage(filter primitive.M) (*int64, error) {
	mg := database.GetDB()
	result, err := mg.MessageCollection.DeleteOne(database.GetDB().MongoCtx, filter)
	if err != nil {
		return nil, err
	}
	if &result != nil {
		return &result.DeletedCount, nil
	}
	return nil, errors.New(errorCodes.NO_RECORD_FOUND)
}

func (m Message) DeleteManyMessages(filter primitive.M) (*int64, error) {
	mg := database.GetDB()
	result, err := mg.MessageCollection.DeleteMany(database.GetDB().MongoCtx, filter)
	if err != nil {
		return nil, err
	}
	if &result != nil {
		return &result.DeletedCount, nil
	}
	return nil, errors.New(errorCodes.NO_RECORD_FOUND)
}

func (m Message) FindOneMessage(filter primitive.M) (*Message, error) {
	mg := database.GetDB()
	var result Message
	err := mg.MessageCollection.FindOne(database.GetDB().MongoCtx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	if &result != nil {
		return &result, nil
	}
	return nil, errors.New(errorCodes.NO_RECORD_FOUND)
}

func (m Message) FindManyMessages(filter primitive.M, findOptions *options.FindOptions) (*[]Message, error) {
	mg := database.GetDB()
	cur, err := mg.MessageCollection.Find(database.GetDB().MongoCtx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(database.GetDB().MongoCtx)
	var messages []Message
	for cur.Next(database.GetDB().MongoCtx) {
		var message Message
		err := cur.Decode(&message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	if &messages != nil {
		return &messages, nil
	}
	return nil, errors.New(errorCodes.NO_RECORD_FOUND)
}

func (m Message) FindOneAndUpdate(filter primitive.M, update primitive.M, opts *options.FindOneAndUpdateOptions) (*Message, error) {
	mg := database.GetDB()
	var message Message
	result := mg.MessageCollection.FindOneAndUpdate(database.GetDB().MongoCtx, filter, update, opts)
	err := result.Decode(&message)
	if err != nil {
		return nil, err
	}
	if &message != nil {
		return &message, nil
	}
	return nil, errors.New(errorCodes.NO_RECORD_FOUND)
}

func (m Message) Aggregate(pipeline mongo.Pipeline) (*[]bson.M, error) {
	mg := database.GetDB()
	cursor, err := mg.MessageCollection.Aggregate(database.GetDB().MongoCtx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(database.GetDB().MongoCtx)
	var results []bson.M
	err = cursor.All(database.GetDB().MongoCtx, &results)
	if err != nil {
		return nil, err
	}
	return &results, nil
}

func PossibleToUpdateStatusWithFilter(nextStatus constants.StatusType, findFilter primitive.M, msgId string, source string) bool {
	var messageModel Message
	msg, err := messageModel.FindOneMessage(findFilter)
	if err != nil {
		logger.LogError(
			nil,
			source+": Status update not possible for msgId: "+msgId+" to -> "+string(nextStatus)+", err: "+err.Error(),
			nil,
		)
		return false
	}
	if msg != nil {
		nextPossibleStates := stateMachine[msg.Status]
		for _, val := range nextPossibleStates {
			if val == nextStatus {
				return true
			}
		}
		logger.LogError(
			nil,
			source+": Status update not possible for msgId: "+msgId+" from "+string(msg.Status)+" to "+string(nextStatus)+" with id: "+msg.Id.Hex(),
			nil,
		)
	}
	return false
}
