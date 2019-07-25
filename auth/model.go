package auth

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type regUser struct {
	ID       *primitive.ObjectID `json:"id" bson:"_id,omitempty" validate:"-"`
	Name     string              `json:"name" bson:"name" validate:"required"`
	Email    string              `json:"email" bson:"email" validate:"required,email"`
	Password string              `json:"password,omitempty" bson:"password" validate:"required"`
}

type loginUser struct {
	ID       *primitive.ObjectID `json:"id" bson:"_id" validate:"-"`
	Name     string              `json:"name" bson:"name"`
	Email    string              `json:"email" bson:"email" validate:"required,email"`
	Password string              `json:"password,omitempty" bson:"password" validate:"required"`
}

func createUser(ctx context.Context, client *mongo.Client, user *regUser) error {
	col := client.Database("codev").Collection("users")

	ok, err := checkExist(ctx, col, "email", user.Email)
	if err != nil {
		return err
	}
	if ok {
		return errUserExists
	}
	ok, err = checkExist(ctx, col, "name", user.Name)
	if err != nil {
		return err
	}
	if ok {
		return errUserExists
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
	fmt.Println(user)
	fmt.Println(insertRes)
	id, ok := insertRes.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("cannot assert id type")
	}
	user.ID = &id
	user.Password = ""

	return nil
}

func checkExist(ctx context.Context, collection *mongo.Collection, key string, value interface{}) (bool, error) {
	findRes := collection.FindOne(ctx, bson.D{
		{key, value},
	})
	err := findRes.Err()
	if err == nil {
		var test interface{}
		err := findRes.Decode(&test)
		if err == nil {
			return true, nil
		} else if err == mongo.ErrNoDocuments {
			return false, nil
		} else {
			return false, err
		}
	} else if err == mongo.ErrNoDocuments {
		return false, nil
	} else {
		return false, err
	}
}

func validateUser(ctx context.Context, client *mongo.Client, user *loginUser) error {
	col := client.Database("codev").Collection("users")
	findRes := col.FindOne(ctx, bson.D{
		{"email", user.Email},
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
