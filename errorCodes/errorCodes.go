package errorCodes

const NO_RECORD_FOUND string = "NO_RECORD_FOUND"

const ERROR_CREATING_NEW_REQUEST string = "ERROR_CREATING_NEW_REQUEST"

type ErrorCodeMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ======== Server Errors ======== //

var BAD_REQUEST = ErrorCodeMessage{
	Code:    1001,
	Message: "Bad Request",
}
var NOT_FOUND = ErrorCodeMessage{
	Code:    1002,
	Message: "Not Found",
}
var INTERNAL_SERVER_ERROR = ErrorCodeMessage{
	Code:    1003,
	Message: "Internal Server Error",
}

// =============================== //

// ======== Database Errors ======== //

var DB_OPERATION_ERROR = ErrorCodeMessage{
	Code:    1050,
	Message: "Database operation error",
}
var DB_OBJECT_ID_CONVERSION_ERROR = ErrorCodeMessage{
	Code:    1051,
	Message: "Error converting to object Id",
}

// =============================== //

// ======== App Errors ======== //

var VALIDATION_ERROR = ErrorCodeMessage{
	Code:    1150,
	Message: "Validation failed",
}
var REQUEST_BODY_PARSING_ERROR = ErrorCodeMessage{
	Code:    1151,
	Message: "Not able to parse request body",
}
var UNMARSHALING_ERROR = ErrorCodeMessage{
	Code:    1152,
	Message: "Unmarshaling error",
}
var MARSHALING_ERROR = ErrorCodeMessage{
	Code:    1153,
	Message: "Marshaling error",
}
var INVALID_TYPE = ErrorCodeMessage{
	Code:    1154,
	Message: "Type is not valid",
}
var INVALID_BODY = ErrorCodeMessage{
	Code:    1155,
	Message: "Body is not valid",
}
var NO_MESSAGE_ID_PROVIDED = ErrorCodeMessage{
	Code:    1156,
	Message: "No message id provided",
}
var BSON_UNMARSHALING_ERROR = ErrorCodeMessage{
	Code:    1157,
	Message: "Unmarshaling error",
}
var BSON_MARSHALING_ERROR = ErrorCodeMessage{
	Code:    1158,
	Message: "Marshaling error",
}
var INVALID_API_KEY_OR_WA_ACCOUNT = ErrorCodeMessage{
	Code:    1159,
	Message: "Not a valid apiKey or wa account",
}
var HEADER_NOT_FOUND = ErrorCodeMessage{
	Code:    1160,
	Message: "x-api-key or x-wa-account header not found",
}
var REQUEST_BODY_TO_BSON_ERROR = ErrorCodeMessage{
	Code:    1161,
	Message: "error converting request body to bson.M",
}
var INVALID_MESSAGE_STATUS = ErrorCodeMessage{
	Code:    1162,
	Message: "message status is not valid",
}
var INTEGER_CONVERSION_ERROR = ErrorCodeMessage{
	Code:    1163,
	Message: "not able to convert to integer",
}
var MESSAGE_LIMIT_EXCEEDED = ErrorCodeMessage{
	Code:    1164,
	Message: "message limit exceeded",
}
var MAKE_MESSAGE_ERROR = ErrorCodeMessage{
	Code:    1165,
	Message: "not able to make message",
}

// =============================== //
