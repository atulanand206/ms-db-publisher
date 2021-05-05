package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/atulanand206/go-kafka"
	"github.com/atulanand206/go-mongo"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	Config struct {
		Row   int    `json:"row" bson:"row, omitempty"`
		Col   int    `json:"col" bson:"col, omitempty"`
		Mines int    `json:"mines" bson:"mines, omitempty"`
		Type  string `json:"name" bson:"name, omitempty"`
	}

	Range struct {
		Start time.Time `json:"start" bson:"start, omitempty"`
		End   time.Time `json:"end" bson:"end, omitempty"`
	}

	Game struct {
		Conf     Config  `json:"config" bson:"config, omitempty"`
		Times    []Range `json:"times" bson:"times, omitempty"`
		Score    int     `json:"score" bson:"score, omitempty"`
		Won      bool    `json:"won" bson:"won, omitempty"`
		Finished bool    `json:"finished" bson:"finished, omitempty"`
	}
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	mongoClientId := os.Getenv("MONGO_CLIENT_ID")
	database := os.Getenv("DATABASE")
	collection := os.Getenv("MONGO_COLLECTION")

	mongo.ConfigureMongoClient(mongoClientId)

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	kafkaBrokerId := os.Getenv("KAFKA_BROKER_ID")

	kafka.LoadConsumer(kafkaBrokerId, kafkaTopic)
	kafka.Read(func(val string) {
		game := Game{}
		json.Unmarshal([]byte(val), &game)
		fmt.Println(game)
		document, _ := document(&game)
		response := mongo.Write(database, collection, *document)
		fmt.Println(response)
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
