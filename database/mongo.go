package database

import (
	"context"
	"fmt"
	"multiBot/config"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoInstance struct {
	Client            *mongo.Client
	Db                *mongo.Database
	MessageCollection *mongo.Collection
	MongoCtx          context.Context
}

var mg MongoInstance

func GetDB() MongoInstance {
	return mg
}

func makeMongoUri() (dbName string, mongoUri string, collectionName string) {

	var mongoUser, mongoUrl, mongoPassword, env string

	cfg := config.GetConfig()

	env = os.Getenv("GO_ENV")
	if env == "" {
		env = "build"
	}

	collectionName = cfg.MONGO.COLLECTION_NAME
	dbName = cfg.MONGO.DB_NAME
	configStr := cfg.MONGO.CONFIG_STR

	mongoUser = os.Getenv("mongoUser")

	mongoPassword = os.Getenv("mongoPassword")

	mongoUrl = cfg.MONGO.URL

	mongoUri = "mongodb://" + mongoUser + ":" + mongoPassword + "@" + mongoUrl + "/" + dbName + configStr
	return

}

func Connect() error {

	allTimeOut := 30 * time.Second

	dbName, mongoUri, collectionName := makeMongoUri()

	clientOptions := options.Client().ApplyURI(mongoUri).SetConnectTimeout(allTimeOut).SetTimeout(allTimeOut)

	ctx := context.Background()

	fmt.Println("Mongo Connecting !")

	//connect to mongoDb
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	fmt.Println("Checking Ping")

	// Checking Ping
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	fmt.Println("Pinged Successful !")

	fmt.Println("Mongo Connected Successfully !")

	dbRef := client.Database(dbName)

	collectionRef := client.Database(dbName).Collection(collectionName)

	mg = MongoInstance{
		Client:            client,
		Db:                dbRef,
		MessageCollection: collectionRef,
		MongoCtx:          ctx,
	}

	return nil

}
