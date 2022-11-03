package main

import (
	// "github.com/pig858/MaxGoLineMongo/db"
	"MaxGoLineMongo/db"
	"context"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spf13/viper"
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

	err = client.Ping(context.TODO(), nil)
	//db test start
	database := client.Database("line")
	messageCollection := database.Collection("Message")
	// testMessage := db.Message{
	// 	Name:    "test",
	// 	Content: "test",
	// 	Time:    time.Now().Unix(),
	// }
	// db.Insert(messageCollection, testMessage)

	if err != nil {
		log.Fatal(err)
	}
	//db connect end
	// targetMessage := db.GetByName(messageCollection, "test")
	// fmt.Println(targetMessage)
	//db test end

	// gin setup test
	r := gin.Default()

	//line bot get message and reply the same one (simple test)
	bot, err := linebot.New(
		viper.GetString("Line.Secret"),
		viper.GetString("Line.Token"),
	)
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/callback", func(c *gin.Context) {
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
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
