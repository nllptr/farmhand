package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User maps database documents to go structs
type User struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Subject string             `bson:"subject,omitempty"`
	Name    string             `bson:"name,omitempty"`
	Email   string             `bson:"email,omitempty"`
}
