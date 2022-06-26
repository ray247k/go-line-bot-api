package util

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mgoCli *mongo.Client

func initEngine() {
	var err error

	// Get configs from file
	viper.SetConfigFile("configs/config.yaml")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig()failed,err:%v\n", err)
		return
	}

	// Set client options
	credential := options.Credential{
		Username: viper.GetString("database.username"),
		Password: viper.GetString("database.password"),
	}
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(credential)

	// Connect to MongoDB
	mgoCli, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	// Check the connection
	err = mgoCli.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
}

func GetMgoCli() *mongo.Client {
	if mgoCli == nil {
		initEngine()
	}
	return mgoCli
}
