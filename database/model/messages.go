package model

import (
	"context"
	"encoding/json"
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
	messagesCollection = client.Database(viper.GetString("database.dbname")).Collection("messages")

	// prepare message
	newMessage := bson.D{
		{"Message", message},
		{"LineId", lineId},
		{"DisplayName", displayName},
	}

	// Insert message into database
	if _, err = messagesCollection.InsertOne(context.TODO(), newMessage); err != nil {
		panic(err.Error())
	}
}

func (Message) GetUserMessageList(lineId string) string {
	filter := bson.M{"LineId": lineId}

	client := util.GetMgoCli()
	messagesCollection = client.Database(viper.GetString("database.dbname")).Collection("messages")

	// get record from database
	cursor, err := messagesCollection.Find(context.TODO(), filter)
	if err != nil {
		panic(err.Error())
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	// fotmat
	messages := make(map[int]string)
	for key, result := range results {
		message := result["Message"]

		if str, ok := message.(string); ok {
			messages[key] = str
		}
	}
	jsonData, _ := json.Marshal(messages)
	return string(jsonData)
}
