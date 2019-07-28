package users

import (
	"context"
	"fmt"

	"github.com/misgorod/co-dev/users/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID    primitive.ObjectID `json:"id" bson:"_id"`
	Name  string             `json:"name,omitempty" bson:"name,omitempty"`
	Email string             `json:"email" bson:"email"`
}

func GetUser(ctx context.Context, client *mongo.Client, id string) (*User, error) {
	coll := client.Database("codev").Collection("users")
	objID, err := primitive.ObjectIDFromHex(id)
	fmt.Println(err)
	if err != nil {
		return nil, errors.ErrUserNotExists
	}
	singleRes := coll.FindOne(ctx, bson.D{
		{
			Key:   "_id",
			Value: objID,
		},
	})
	var user User
	err = singleRes.Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrUserNotExists
		}
		return nil, err
	}
	return &user, nil
}
