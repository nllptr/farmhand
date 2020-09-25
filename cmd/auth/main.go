package main

import (
	"context"
	"net/http"
	"os"

	"github.com/coreos/go-oidc"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/nllptr/farmhand/pkg/auth"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// Required environment variables
const (
	AuthClientID     = "AUTH_CLIENT_ID"
	AuthClientSecret = "AUTH_CLIENT_SECRET"
	AuthRedirectURL  = "AUTH_REDIRECT_URL"
	MongoDbURI       = "MONGODB_URI"
	MongoDbName      = "MONGODB_NAME"
)

func main() {
	logger, _ := zap.NewDevelopment()
	slog := logger.Sugar()

	err := godotenv.Load()
	if err != nil {
		_, ok := os.LookupEnv(AuthClientID)
		if !ok {
			slog.Fatal("environment variable AUTH_CLIENT_ID not set")
		}
		_, ok = os.LookupEnv(AuthClientSecret)
		if !ok {
			slog.Fatal("environment variable AUTH_CLIENT_SECRET not set")
		}
		_, ok = os.LookupEnv(AuthRedirectURL)
		if !ok {
			slog.Fatal("environment variable AUTH_RECIRECT_URL not set")
		}
		_, ok = os.LookupEnv(MongoDbURI)
		if !ok {
			slog.Fatal("environment variable MONGODB_URI not set")
		}
		_, ok = os.LookupEnv(MongoDbName)
		if !ok {
			slog.Fatal("environment variable MONGODB_NAME not set")
		}
	}
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		slog.Fatalw("failed to create OIDC provider: ", err)
	}
	config := &oauth2.Config{
		ClientID:     os.Getenv(AuthClientID),
		ClientSecret: os.Getenv(AuthClientSecret),
		RedirectURL:  os.Getenv(AuthRedirectURL),
		Endpoint:     provider.Endpoint(),
		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv(MongoDbURI)))
	if err != nil {
		slog.Fatalw("failed to create mongodb client:", err)
	}
	err = client.Connect(ctx)
	if err != nil {
		slog.Fatalw("database connection failed:", err)
	}
	defer client.Disconnect(ctx)
	db := client.Database(os.Getenv(MongoDbName))

	r := mux.NewRouter()
	r.HandleFunc("/auth/google", auth.CreateRedirect(config))
	r.HandleFunc("/auth/callback", auth.CreateCallback(provider, config, db))

	slog.Info("Auth is looking good...")
	slog.Fatal(http.ListenAndServe(":8080", r))
}
