package models

import (
	"context"
	"github.com/misgorod/co-dev/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"io"
)

type File struct {
	ID primitive.ObjectID `json:"id"`
}

func DownloadFile(ctx context.Context, client *mongo.Client, id string, writer io.Writer) error {
	db := client.Database("codev")
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		return err
	}
	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = bucket.DownloadToStream(obj, writer)
	if err != nil {
		if err == gridfs.ErrFileNotFound {
			return errors.ErrNoFile
		}
		return err
	}
	return nil
}
