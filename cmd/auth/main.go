package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/coreos/go-oidc"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/nllptr/farmhand/pkg/auth"
	"golang.org/x/oauth2"
)

// AuthClientID, AuthClientSecret and Host are required environment variables.
const (
	AuthClientID     = "AUTH_CLIENT_ID"
	AuthClientSecret = "AUTH_CLIENT_SECRET"
	Host             = "HOST"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		_, ok := os.LookupEnv(AuthClientID)
		if !ok {
			log.Fatal("environment variable AUTH_CLIENT_ID not set")
		}
		_, ok = os.LookupEnv(AuthClientSecret)
		if !ok {
			log.Fatal("environment variable AUTH_CLIENT_SECRET not set")
		}
		_, ok = os.LookupEnv(Host)
		if !ok {
			log.Fatal("environment variable HOST not set")
		}
	}
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		log.Panic("failed to create OIDC provider: ", err)
	}
	config := &oauth2.Config{
		ClientID:     os.Getenv(AuthClientID),
		ClientSecret: os.Getenv(AuthClientSecret),
		RedirectURL:  os.Getenv(Host) + "/auth/callback",
		Endpoint:     provider.Endpoint(),
		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}
	r := mux.NewRouter()
	r.HandleFunc("/auth/google", auth.CreateRedirect(config))
	r.HandleFunc("/auth/callback", auth.CreateCallback(provider, config))

	log.Fatal(http.ListenAndServe(":8080", r))
}
