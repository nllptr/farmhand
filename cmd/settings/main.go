package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"github.com/nllptr/farmhand/pkg/settings"
)

// GCPProjectID environment variable constant
const GCPProjectID = "GCP_PROJECT_ID"

func main() {
	_, ok := os.LookupEnv(GCPProjectID)
	if !ok {
		log.Fatal("environment variable GCP_PROJECT_ID not set")
	}
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv(GCPProjectID))
	defer client.Close()
	if err != nil {
		log.Fatal("failed to create Firestore client:", err)
	}

	r := mux.NewRouter()

	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/settings", settings.CreateGetSettings(client))

	log.Fatal(http.ListenAndServe(":8080", r))
}
