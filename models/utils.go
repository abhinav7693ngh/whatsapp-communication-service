package models

import "go.mongodb.org/mongo-driver/bson/primitive"

func GetStringObjectIdFromBson(objId interface{}) string {
	return objId.(primitive.ObjectID).Hex()
}
