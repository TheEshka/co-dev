package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID       *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Email    string              `json:"email" bson:"email"`
	Password string              `json:"password,omitempty" bson:"password,omitempty"`
}

func createUser(ctx context.Context, client *mongo.Client, email string, password string) (*User, error) {
	col := client.Database("codev").Collection("users")

	findRes := col.FindOne(ctx, bson.D{
		{"email", email},
	})
	if findRes.Err() == nil {
		var testUser User
		err := findRes.Decode(&testUser)
		if err != nil && err != mongo.ErrNoDocuments {
			return nil, err
		}
		if err == nil {
			return nil, ErrUserExists
		}
	}
	if err := findRes.Err(); err != mongo.ErrNoDocuments && err != nil {
		return nil, findRes.Err()
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:    email,
		Password: string(hash),
	}

	insertRes, err := col.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	id, ok := insertRes.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("cannot assert id type")
	}
	user.ID = &id
	user.Password = ""

	return user, nil
}

func validateUser(ctx context.Context, client *mongo.Client, email string, password string) (*User, error) {
	col := client.Database("codev").Collection("users")
	var user User
	findRes := col.FindOne(ctx, bson.D{
		{"email", email},
	})
	err := findRes.Decode(&user)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return &user, nil
}
