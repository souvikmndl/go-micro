package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/souvikmndl/logger/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	grpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongo
	mongoCient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoCient
	log.Println("connected to mongo")

	// create a context to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
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

	log.Println("starting server on port:", webPort)
	err = srv.ListenAndServe()
	if err != nil {
		log.Panicf("unable to start logger service %s\n", err.Error())
	}
	log.Println("server started on port ", webPort)
}

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panicf("unable to start logger service %s\n", err.Error())
	}
}

func connectToMongo() (*mongo.Client, error) {
	// create the connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("error connecting to mongo: ", err)
		return nil, err
	}

	return c, nil
}
