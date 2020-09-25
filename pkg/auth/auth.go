package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/nllptr/farmhand/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
)

// TODO: This entire file needs an overhaul once work on the frontend has started.

// CreateRedirect returns a handler function that redirects to OpenID Connect
func CreateRedirect(c *oauth2.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// TODO: The state variable should come from the browser client.
		// Should the final check also somehow be in the browser?

		state := randStateString()
		http.SetCookie(w, &http.Cookie{
			Name:  "state",
			Value: state,
		})

		http.Redirect(w, r, c.AuthCodeURL(state), http.StatusFound)
	}
}

// CreateCallback returns a handler function that verifies an OIDC token.
func CreateCallback(p *oidc.Provider, c *oauth2.Config, d *mongo.Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		stateCookie, err := r.Cookie("state")
		if err != nil {
			http.Error(w, "No state cookie present.", http.StatusBadRequest)
			return
		}
		if r.URL.Query().Get("state") != stateCookie.Value {
			http.Error(w, "State did not match.", http.StatusBadRequest)
			return
		}
		// state checked, clean up temporary cookie
		http.SetCookie(w, &http.Cookie{
			Name:   "state",
			MaxAge: -1,
		})

		oauth2Token, err := c.Exchange(r.Context(), r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
			return
		}
		oidcConfig := &oidc.Config{
			ClientID: c.ClientID,
		}
		idToken, err := p.Verifier(oidcConfig).Verify(r.Context(), rawIDToken)
		if err != nil {
			http.Error(w, "Token validation failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		var claims struct {
			Email   string `json:"email"`
			Name    string `json:"name"`
			Subject string `json:"sub"`
		}
		err = idToken.Claims(&claims)
		if err != nil {
			http.Error(w, "Claims retreival failed:"+err.Error(), http.StatusInternalServerError)
		}

		user := db.User{}

		err = d.Collection("users").FindOne(r.Context(), bson.M{"sub": claims.Subject}).Decode(&user)
		if err != mongo.ErrNilDocument {
			// register
			http.Error(w, "find user:"+err.Error(), http.StatusInternalServerError)
		}
		if err != nil {
			// other error
		}

		fmt.Fprintf(w, "%v", claims)
		fmt.Fprintf(w, "rawIDToken: %v", rawIDToken)
	}
}

func randStateString() string {
	var state strings.Builder
	charSet := "ABCDEDFGHIJKLMNOPQRSTabcdedfghijklmnopqrst"
	length := 30
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		state.WriteString(string(randomChar))
	}
	b64State := base64.URLEncoding.EncodeToString([]byte(state.String()))

	return b64State
}

type key int

const (
	// KeyUserID is used for adding/retriving the user ID to/from the request context.
	KeyUserID key = iota
)

// CreateVerifyJWTMiddleware verifies that the JWT is present and valid, then adds the subject field to the request context.
func CreateVerifyJWTMiddleware(next http.Handler, p *oidc.Provider, c *oidc.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		split := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(split) != 2 {
			http.Error(w, "Malformed Authorization header", http.StatusUnauthorized)
			return
		}
		idToken, err := p.Verifier(c).Verify(r.Context(), split[1])
		if err != nil {
			http.Error(w, "Invalid access token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), KeyUserID, idToken.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
