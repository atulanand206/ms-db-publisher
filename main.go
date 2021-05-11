package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/atulanand206/go-kafka"
	"github.com/atulanand206/go-mongo"
	"github.com/atulanand206/ms-db-publisher/objects"
	"github.com/atulanand206/ms-db-publisher/routes"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	godotenv.Load()

	mongoClientId := os.Getenv("MONGO_CLIENT_ID")
	database := os.Getenv("DATABASE")
	collection := os.Getenv("MONGO_COLLECTION")

	mongo.ConfigureMongoClient(mongoClientId)

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	kafkaBrokerId := os.Getenv("KAFKA_BROKER_ID")

	kafka.LoadConsumer(kafkaBrokerId, kafkaTopic)
	kafka.Read(func(val string) {
		game := objects.Game{}
		json.Unmarshal([]byte(val), &game)
		fmt.Println(game)
		document, _ := document(&game)
		response := mongo.Write(database, collection, *document)
		fmt.Println(response)
		user, err := routes.GetUser(game.Player.Username)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(user)
		var userRequest objects.UserRequest
		userRequest.Username = user.Username
		userRequest.Name = user.Name
		userRequest.Rating = user.Rating + game.Score
		fmt.Println(userRequest)
		res, err := routes.UpdateUser(user.Id, userRequest)
		fmt.Println(res)
		if err != nil {
			fmt.Println(err)
			return
		}
	})
}

func document(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		log.Panic(err)
	}
	err = bson.Unmarshal(data, &doc)
	return
}
