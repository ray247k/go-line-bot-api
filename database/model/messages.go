package model

import (
	"context"
	"fmt"
	"line-webhook-receiver/database/util"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Message struct {
	Message     string `bson:"message"`
	LineId      string `bson:"lineId"`
	DisplayName string `bson:"displayName"`
}

var (
	messagesCollection *mongo.Collection
	err                error
)

func (Message) Store(message string, lineId string, displayName string) {
	// Connect to mongoDB
	client := util.GetMgoCli()
	defer client.Disconnect(context.TODO())
	messagesCollection = client.Database(viper.GetString("database.dbname")).Collection("messages")

	// prepare message
	newMessage := bson.D{
		{"Message", message},
		{"LineId", lineId},
		{"DisplayName", displayName},
	}

	// Insert message into database
	if _, err = messagesCollection.InsertOne(context.TODO(), newMessage); err != nil {
		fmt.Print(err)
		return
	}
}
