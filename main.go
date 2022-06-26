package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	model "line-webhook-receiver/database/model"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spf13/viper"
)

var (
	lineBotClient *linebot.Client
	lineBotErr    error
)

func main() {
	// Get configs from file
	viper.SetConfigFile("configs/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig()failed,err:%v\n", err)
		return
	}

	// Init line bot
	lineBotChannelSecret, lineBotChannelAccessToken := viper.GetString("line-sdk.channel-secret"), viper.GetString("line-sdk.channel-access-token")
	lineBotClient, lineBotErr = linebot.New(lineBotChannelSecret, lineBotChannelAccessToken)
	if lineBotErr != nil {
		panic(err.Error())
	}

	router := gin.Default()
	messages := router.Group("/messages")
	{
		messages.POST("/callback", storeMessage)
		messages.POST("/send", sendMessage)
		messages.GET("/:userId", userMessages)
	}

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// Receive message from line webhook, save the user info and message in MongoDB
func storeMessage(c *gin.Context) {
	events, err := lineBotClient.ParseRequest(c.Request)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Bad Request",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Bad Request",
			})
		}
		return
	}

	newMessage := ""
	lineId := ""
	displayName := ""

	for _, event := range events {
		lineId = event.Source.UserID
		profile, err := lineBotClient.GetProfile(event.Source.UserID).Do()
		if err != nil {
			panic(err.Error())
		}
		displayName = profile.DisplayName

		fmt.Println(profile)
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				newMessage = message.Text
			}
		}
	}

	message := new(model.Message)
	message.Store(newMessage, lineId, displayName)

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
}

// Send message back to line
func sendMessage(c *gin.Context) {
	message := c.PostForm("message")
	targetUserId := viper.GetString("line.target-user-id")

	if _, err := lineBotClient.PushMessage(targetUserId, linebot.NewTextMessage(message)).Do(); err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
}

// Query message list of the user from MongoDB
func userMessages(c *gin.Context) {
	userId := c.Param("userId")

	message := new(model.Message)
	messages := message.GetUserMessageList(userId)

	c.JSON(http.StatusOK, gin.H{
		"messages": json.RawMessage(messages),
	})
}
