package models

import (
	"context"
	errors2 "github.com/misgorod/co-dev/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type RegUser struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty" validate:"-"`
	Name     string             `json:"name" bson:"name" validate:"required"`
	Email    string             `json:"email" bson:"email" validate:"required,email"`
	Password string             `json:"password,omitempty" bson:"password" validate:"required"`
}

type LoginUser struct {
	ID       primitive.ObjectID `json:"id" bson:"_id" validate:"-"`
	Name     string             `json:"name" bson:"name"`
	Email    string             `json:"email" bson:"email" validate:"required,email"`
	Password string             `json:"password,omitempty" bson:"password" validate:"required"`
}

func CreateUser(ctx context.Context, client *mongo.Client, user *RegUser) error {
	col := client.Database("codev").Collection("users")

	ok, err := CheckExist(ctx, col, "email", user.Email)
	if err != nil {
		return err
	}
	if ok {
		return errors2.ErrUserExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)

	insertRes, err := col.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	id, ok := insertRes.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors2.ErrAssertID
	}
	user.ID = id
	user.Password = ""

	return nil
}

func ValidateUser(ctx context.Context, client *mongo.Client, user *LoginUser) error {
	col := client.Database("codev").Collection("users")
	findRes := col.FindOne(ctx, bson.D{
		{
			Key:   "email",
			Value: user.Email,
		},
	})
	var dbUser LoginUser
	err := findRes.Decode(&dbUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors2.ErrWrongCreds
		}
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return errors2.ErrWrongCreds
		}
		return err
	}
	*user = dbUser
	user.Password = ""
	return nil
}
