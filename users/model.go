package users

import (
	"context"
	"fmt"

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
		return nil, ErrUserNotExists
	}
	singleRes := coll.FindOne(ctx, bson.D{
		{
			Key:   "_id",
			Value: objID,
		},
	})
	fmt.Println(singleRes)
	fmt.Println(singleRes.Err())

	if singleRes.Err() != nil {
		if singleRes.Err() == mongo.ErrNoDocuments {
			return nil, ErrUserNotExists
		}
		return nil, singleRes.Err()
	}
	var user User
	err = singleRes.Decode(&user)
	fmt.Println(err)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotExists
		}
		return nil, err
	}
	return &user, nil
}
