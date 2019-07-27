package common

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckExist(ctx context.Context, collection *mongo.Collection, key string, value interface{}) (bool, error) {
	findRes := collection.FindOne(ctx, bson.D{
		{
			Key:   key,
			Value: value,
		},
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
