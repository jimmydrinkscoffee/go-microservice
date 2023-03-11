package main

import (
	"context"
	"fmt"
	"log"
	"logger/data"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoUri = "mongodb://mongo:27017"
	grpcPort = "50001"
)

var client *mongo.Client

type Config struct{ Models data.Models }

func main() {
	// connect to Mongo

	mgClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mgClient

	// create context to disconnect to Mongo

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	if err = srv.ListenAndServe(); err != nil {
		log.Panicln("Error server: ", err)
	} else {
		log.Println("Server listen to: ", webPort)
	}
}

func connectToMongo() (*mongo.Client, error) {
	clOpts := options.Client().ApplyURI(mongoUri)
	clOpts.SetAuth(options.Credential{
		Username: "admin",
		Password: "mongo_password",
	})

	conn, err := mongo.Connect(context.TODO(), clOpts)
	if err != nil {
		log.Println("Error conn: ", err)
		return nil, err
	}

	return conn, nil
}
