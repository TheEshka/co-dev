package models

import (
	"context"
	errors2 "github.com/misgorod/co-dev/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"io"
)

type User struct {
	ID    primitive.ObjectID `json:"id" bson:"_id"`
	Name  string             `json:"name,omitempty" bson:"name,omitempty"`
	Email string             `json:"email" bson:"email"`
	Info  UserInfo           `json:"info,omitempty" bson:"info,omitempty"`
}

type UserInfo struct {
	ImageID *primitive.ObjectID `json:"image,omitempty" bson:"image,omitempty"`
	City string                 `json:"city" bson:"city"`
	GithubLink string           `json:"githubLink" bson:"githubLink"`
	AboutMe string              `json:"aboutMe" bson:"aboutMe"`
	AuthorPosts []*Post         `json:"authorPosts" bson:"authorPosts"`
	MemberPosts []*Post         `json:"memberPosts" bson:"memberPosts"`
}

func GetUser(ctx context.Context, client *mongo.Client, id string) (*User, error) {
	users := client.Database("codev").Collection("users")
	posts := client.Database("codev").Collection("posts")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors2.ErrUserNotExists
	}
	singleRes := users.FindOne(ctx, bson.D{
		{
			Key:   "_id",
			Value: objID,
		},
	})
	var user User
	err = singleRes.Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors2.ErrUserNotExists
		}
		return nil, err
	}
	filterAuthor := bson.D{{"author._id", objID}}
	curAuthor, err := posts.Find(ctx, filterAuthor)
	if err != nil {
		return nil, err
	}
	defer curAuthor.Close(ctx)
	user.Info.AuthorPosts = make([]*Post, 0)
	for curAuthor.Next(ctx) {
		var post Post
		err := curAuthor.Decode(&post)
		if err != nil {
			return nil, err
		}
		user.Info.AuthorPosts = append(user.Info.AuthorPosts, &post)
	}

	filterMember := bson.D{{"members._id", objID}}
	curMember, err := posts.Find(ctx, filterMember)
	if err != nil {
		return nil, err
	}
	defer curMember.Close(ctx)
	user.Info.MemberPosts = make([]*Post, 0)
	for curMember.Next(ctx) {
		var post Post
		err := curMember.Decode(&post)
		if err != nil {
			return nil, err
		}
		user.Info.MemberPosts = append(user.Info.MemberPosts, &post)
	}
	return &user, nil
}

func PutUser(ctx context.Context, client *mongo.Client, userID string, info *UserInfo) (*User, error) {
	coll := client.Database("codev").Collection("users")
	user, err := GetUser(ctx, client, userID)
	if err != nil {
		return nil, err
	}
	if info.AboutMe != "" {
		user.Info.AboutMe = info.AboutMe
	}
	if info.City != "" {
		user.Info.City = info.City
	}
	if info.GithubLink != "" {
		user.Info.GithubLink = info.GithubLink
	}
	filter := bson.D{{"_id", userID}}
	update := bson.D{{"$set", bson.D{{"info", user.Info}}}}
	_, err = coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func UploadUserImage(ctx context.Context, client *mongo.Client, reader io.Reader, user *User) (*File, error) {
	db := client.Database("codev")
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		return nil, err
	}
	fileID, err := bucket.UploadFromStream("image", reader)
	if err != nil {
		return nil, err
	}
	col := db.Collection("posts")
	filter := bson.D{{"_id", user.ID}}
	update := bson.D{{"$set", bson.D{
		{"info.image", fileID},
	}}}
	_, err = col.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	user.Info.ImageID = &fileID

	return &File{fileID}, nil
}