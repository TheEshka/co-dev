package auth

import (
	"context"

	"github.com/misgorod/co-dev/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type regUser struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty" validate:"-"`
	Name     string             `json:"name" bson:"name" validate:"required"`
	Email    string             `json:"email" bson:"email" validate:"required,email"`
	Password string             `json:"password,omitempty" bson:"password" validate:"required"`
}

type loginUser struct {
	ID       primitive.ObjectID `json:"id" bson:"_id" validate:"-"`
	Name     string             `json:"name" bson:"name"`
	Email    string             `json:"email" bson:"email" validate:"required,email"`
	Password string             `json:"password,omitempty" bson:"password" validate:"required"`
}

func createUser(ctx context.Context, client *mongo.Client, user *regUser) error {
	col := client.Database("codev").Collection("users")

	ok, err := common.CheckExist(ctx, col, "email", user.Email)
	if err != nil {
		return err
	}
	if ok {
		return ErrUserExists
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
		return ErrAssertID
	}
	user.ID = id
	user.Password = ""

	return nil
}

func validateUser(ctx context.Context, client *mongo.Client, user *loginUser) error {
	col := client.Database("codev").Collection("users")
	findRes := col.FindOne(ctx, bson.D{
		{
			Key:   "email",
			Value: user.Email,
		},
	})
	var dbUser loginUser
	err := findRes.Decode(&dbUser)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return err
	}
	*user = dbUser
	user.Password = ""
	return nil
}
