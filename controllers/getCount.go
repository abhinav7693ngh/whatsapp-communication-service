package controllers

import (
	"multiBot/errorCodes"
	"multiBot/models"
	"multiBot/types"
	"multiBot/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ValiadateGetCountPayload(c *fiber.Ctx) error {
	var getCountPayload types.GetCountPayload
	err := c.BodyParser(&getCountPayload)
	if err != nil {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: err.Error(),
				ErrorCode:    errorCodes.REQUEST_BODY_PARSING_ERROR,
			},
			getCountPayload,
		)
	}

	v := validate.Struct(getCountPayload)
	if !v.Validate() {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: v.Errors.One(),
				ErrorCode:    errorCodes.VALIDATION_ERROR,
			},
			getCountPayload,
		)
	}

	c.Locals("getCountPayload", getCountPayload)
	return c.Next()
}

func GetCount(c *fiber.Ctx) error {
	getCountPayload := c.Locals("getCountPayload").(types.GetCountPayload)

	var matchConditions []bson.D
	for _, msgStatusValue := range getCountPayload.Status {
		matchConditions = append(matchConditions, bson.D{{"status", msgStatusValue}})
	}

	matchStage := bson.D{{"$match", bson.D{
		{"clientIdentifier", bson.D{{"$in", getCountPayload.Clients}}},
		{"$or", matchConditions},
		{"$and", bson.A{
			bson.D{{"createdAt", bson.D{{"$gte", getCountPayload.StartTime}}}},
			bson.D{{"updatedAt", bson.D{{"$lt", getCountPayload.EndTime}}}},
		}},
	}}}

	pipeline := mongo.Pipeline{
		matchStage,
		bson.D{{"$group", bson.D{{"_id", "$status"}, {"count", bson.D{{"$sum", 1}}}}}},
	}

	var messageModel models.Message
	results, err := messageModel.Aggregate(pipeline)
	if err != nil {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: err.Error(),
				ErrorCode:    errorCodes.DB_OPERATION_ERROR,
			},
			getCountPayload,
		)
	}

	var counts []bson.M
	for _, result := range *results {
		count := bson.M{
			"count":  result["count"],
			"status": result["_id"],
		}
		counts = append(counts, count)
	}

	return utils.SuccessResponse(
		c,
		counts,
		getCountPayload,
	)
}
