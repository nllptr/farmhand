package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/nllptr/farmhand/pkg/settings"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBURI environment variable constant
const mongodbURI = "MONGODB_URI"

func main() {

	uri, ok := os.LookupEnv(mongodbURI)
	if !ok {
		log.Fatal("environment variable MONGODB_URI not set")
	}
	ctx := context.Background()

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	databases, err := client.ListDatabaseNames(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("databases: %v\n", databases)

	r := mux.NewRouter()

	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/settings", settings.CreateGetSettings(client))

	log.Fatal(http.ListenAndServe(":8080", r))
}
