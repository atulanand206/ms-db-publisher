package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/atulanand206/go-kafka"
	"github.com/atulanand206/go-mongo"
	"github.com/atulanand206/ms-db-publisher/objects"
	"github.com/atulanand206/ms-db-publisher/routes"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// Register the MongoDB cloud atlas database.
	mongoClientId := os.Getenv("MONGO_CLIENT_ID")
	database := os.Getenv("DATABASE")
	collection := os.Getenv("MONGO_COLLECTION")

	mongo.ConfigureMongoClient(mongoClientId)

	// Register the Kafka cluster subscriber.
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	kafkaBrokerId := os.Getenv("KAFKA_BROKER_ID")

	kafka.LoadConsumer(kafkaBrokerId, kafkaTopic)
	kafka.Read(func(val string) {
		// Decode the kafka response from message.
		r := &objects.Response{}
		json.Unmarshal([]byte(val), &r)
		var game = r.Match
		document, _ := mongo.Document(&game)
		// Create a game in the database.
		mongo.Write(database, collection, *document)
		// Get the user from the users service.
		user, err := routes.GetUser(game.Player.Username, r.Token)
		if err != nil {
			return
		}
		// Create the  user request
		var userRequest objects.UserRequest
		userRequest.Username = user.Username
		userRequest.Name = user.Name
		userRequest.Rating = user.Rating + game.Score
		// Update the user in the database.
		_, err = routes.UpdateUser(user.Id, userRequest, r.Token)
		if err != nil {
			fmt.Println(err)
			return
		}
	})
}
