package main

import (
	"MaxGoLineMongo/db"
	"context"
	"fmt"
	"time"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()

	if err != nil {
		panic("Cannot read config file: " + err.Error())
	}

	// db connent start
	client := db.Connect(viper.GetString("DB.Host"), viper.GetString("DB.Port"))
	defer client.Disconnect(context.TODO())
	//db connect end

	//get db collection start
	database := client.Database("line")
	messageCollection := database.Collection("Message")
	// followerCollection := database.Collection("Follower")
	//get db collection end

	// gin setup
	r := gin.Default()

	//init line bot
	bot, err := linebot.New(
		viper.GetString("Line.Secret"),
		viper.GetString("Line.Token"),
	)
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/receive", func(c *gin.Context) {
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				log.Print(err)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					res, err := bot.GetProfile(event.Source.UserID).Do()
					if err != nil {
						log.Print(err)
					}

					receiveMessage := db.Message{
						UserID:  event.Source.UserID,
						Name:    res.DisplayName,
						Content: message.Text,
						Time:    time.Now().Format("2006-01-02 15:04:05"),
					}
					db.Insert(messageCollection, receiveMessage)
				}
			}
		}
	})

	r.GET("/sendMsgToLine", func(c *gin.Context) {
		if _, err := bot.PushMessage("U24407c0d7e4a7d85dddc3140cff3d6ea", linebot.NewTextMessage("hello")).Do(); err != nil {
			panic(err)
		}
	})

	r.GET("/queryMsg", func(c *gin.Context) {
		filter := bson.M{}
		datas := db.GetByFilter[db.Message](messageCollection, filter)

		for _, v := range datas {
			fmt.Printf("name : %s\n", v.Name)
			fmt.Printf("content : %s\n", v.Content)
			fmt.Printf("time : %s\n", v.Time)
		}
	})

	r.GET("/queryMsg/:name", func(c *gin.Context) {
		filter := bson.D{{"name", "楷鎰"}}
		datas := db.GetByFilter[db.Message](messageCollection, filter)

		for _, v := range datas {
			fmt.Printf("name : %s\n", v.Name)
			fmt.Printf("content : %s\n", v.Content)
			fmt.Printf("time : %s\n", v.Time)
		}
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
