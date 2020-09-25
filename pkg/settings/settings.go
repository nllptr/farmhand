package settings

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/nllptr/farmhand/pkg/auth"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserSettings holds the settings for a user
type UserSettings struct {
	Name               string `firestore:"name" json:"name"`
	Email              string `firestore:"email" json:"email"`
	Phone              string `firestore:"phone" json:"phone"`
	TempActiveEmail    bool   `firestore:"tempActiveEmail" json:"tempActiveEmail"`
	TempActivePhone    bool   `firestore:"tempActivePhone" json:"tempActivePhone"`
	TempActiveSchedule string `firestore:"tempActiveSchedule" json:"tempActiveSchedule"`
	DiseaseActiveEmail bool   `firestore:"diseaseActiveEmail" json:"diseaseActiveEmail"`
	DiseaseActivePhone bool   `firestore:"diseaseActivePhone" json:"diseaseActivePhone"`
}

// CreateGetSettings creates a GetSettings handler
func CreateGetSettings(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sub := getUserID(w, r)
		fmt.Fprintf(w, "sub: %v\n", sub)

		// dsnap, err := client.Collection("users").Doc(sub).Get(r.Context())
		// if err != nil {
		// 	http.Error(w, "Firestore error", http.StatusInternalServerError)
		// 	fmt.Printf("Firestore error: %v", err)
		// 	return
		// }
		// var settings UserSettings
		// dsnap.DataTo(&settings)
		// settings.Username = dsnap.Ref.ID

		w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(settings)
	}
}

// CreatePatchSettings creates a PatchSettings handler
func CreatePatchSettings(client *firestore.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Malformed request body", http.StatusBadRequest)
			fmt.Printf("Malformed request body: %v", err)
			return
		}

		keys := []string{
			"name",
			"email",
			"phone",
			"tempActiveEmail",
			"tempActivePhone",
			"tempActiveSchedule",
			"diseaseActiveEmail",
			"diseaseActivePhon",
		}
		sub := getUserID(w, r)
		var updates []firestore.Update

		for _, key := range keys {
			newValue, isUpdated := body[key]
			if isUpdated {
				updates = append(updates, firestore.Update{Path: key, Value: newValue})
			}
		}
		user := client.Collection("users").Doc(sub)
		_, err = user.Update(r.Context(), updates)
		if err != nil {
			http.Error(w, "Firestore update failed", http.StatusInternalServerError)
			fmt.Printf("Firestore update failed: %v", err)
		}
	}
}

func getUserID(w http.ResponseWriter, r *http.Request) string {
	userID, ok := r.Context().Value(auth.KeyUserID).(string)
	if !ok {
		http.Error(w, "No user ID in request.", http.StatusUnauthorized)
		fmt.Printf("No user ID in request")
	}
	return userID
}
