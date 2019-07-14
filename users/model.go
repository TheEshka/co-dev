package users

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
	Name     string              `json:"name,omitempty" bson:"name,omitempty"`
	Email    string              `json:"email" bson:"email"`
	Password string              `json:"password,omitempty" bson:"password,omitempty"`
}

func CreateUser(ctx context.Context, client *mongo.Client, user *User) error {
	col := client.Database("codev").Collection("users")

	ok, err := CheckExist(ctx, col, "email", user.Email)
	if err != nil {
		return err
	}
	if ok {
		return ErrUserExists
	}
	ok, err = CheckExist(ctx, col, "name", user.Name)
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
		return errors.New("cannot assert id type")
	}
	user.ID = &id
	user.Password = ""

	return nil
}

func CheckExist(ctx context.Context, collection *mongo.Collection, key string, value interface{}) (bool, error) {
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

func ValidateUser(ctx context.Context, client *mongo.Client, user *User) error {
	col := client.Database("codev").Collection("users")
	findRes := col.FindOne(ctx, bson.D{
		{"email", user.Email},
	})
	var dbUser User
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

func GetUser(ctx context.Context, client *mongo.Client, id string) (*User, error) {
	coll := client.Database("codev").Collection("users")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrUserNotExists
	}
	singleRes := coll.FindOne(ctx, bson.D{{"_id", objId}})
	if singleRes.Err() != nil {
		if singleRes.Err() == mongo.ErrNoDocuments {
			return nil, ErrUserNotExists
		}
		return nil, singleRes.Err()
	}
	var user User
	err = singleRes.Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotExists
		}
		return nil, err
	}
	user.Password = ""
	return &user, nil
}
