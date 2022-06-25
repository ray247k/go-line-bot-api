package main

import (
	"fmt"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spf13/viper"
)

var (
	lineBotClient *linebot.Client
	lineBotErr    error
)

func main() {
	viper.SetConfigFile("configs/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig()failed,err:%v\n", err)
		return
	}

	lineBotChannelSecret, lineBotChannelAccessToken := viper.GetString("line-sdk.channel-secret"), viper.GetString("line-sdk.channel-access-token")
	lineBotClient, lineBotErr = linebot.New(lineBotChannelSecret, lineBotChannelAccessToken)

}
