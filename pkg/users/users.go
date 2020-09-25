package users

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
}

type MongoService struct {
	db mongo.Database
}

func New(URI string) MongoService {
	// client, err := mongo.NewClient(options.Client().ApplyURI(URI))
	// if err != nil {
	// }
	return MongoService{}
}
