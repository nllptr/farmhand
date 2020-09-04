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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file", err)
	}
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		log.Panic("failed to create OIDC provider: ", err)
	}
	config := &oauth2.Config{
		ClientID:     os.Getenv("AUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("HOST") + "/auth/callback",
		Endpoint:     provider.Endpoint(),
		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}
	r := mux.NewRouter()
	r.HandleFunc("/auth/google", auth.CreateRedirect(config))
	r.HandleFunc("/auth/callback", auth.CreateCallback(provider, config))

	log.Fatal(http.ListenAndServe(":8080", r))
}
