package main

import (
	"context"
	"fmt"

	util "line-webhook-receiver/database/util"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	lineBotClient      *linebot.Client
	lineBotErr         error
	messagesCollection *mongo.Collection
	insertResult       *mongo.InsertOneResult
)

func main() {
	// Get configs from file
	viper.SetConfigFile("configs/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig()failed,err:%v\n", err)
		return
	}

	// Connect to mongoDB
	client := util.GetMgoCli()
	defer client.Disconnect(context.TODO())
	messagesCollection = client.Database(viper.GetString("database.dbname")).Collection("messages")

	// Init line bot
	lineBotChannelSecret, lineBotChannelAccessToken := viper.GetString("line-sdk.channel-secret"), viper.GetString("line-sdk.channel-access-token")
	lineBotClient, lineBotErr = linebot.New(lineBotChannelSecret, lineBotChannelAccessToken)

	// TODO prepare message
	message := bson.D{{"Message", "Hi"}, {"UserInfo", "Nanami"}}

	// Insert message into database
	if insertResult, err = messagesCollection.InsertOne(context.TODO(), message); err != nil {
		fmt.Print(err)
		return
	}
}
