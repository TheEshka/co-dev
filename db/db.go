package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() (*mongo.Client, error) {
	dbOpts := options.Client().ApplyURI("mongodb://mongo:27017")
	fmt.Println("Start connect")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, conErr := mongo.Connect(ctx, dbOpts)

	if conErr != nil {
		fmt.Println(conErr)
		return nil, conErr
	}

	return client, nil
}
