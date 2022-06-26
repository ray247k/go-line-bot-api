package model

type Message struct {
	Message  string `bson:"message"`
	UserInfo string `bson:"userInfo"`
}
