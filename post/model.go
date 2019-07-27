package post

import (
	"context"
	"errors"
	"time"

	"github.com/misgorod/co-dev/auth"
	"github.com/misgorod/co-dev/users"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
	ID             *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title          string              `json:"title" bson:"title" validate:"required,gte=5"`
	Subject        string              `json:"subject" bson:"subject" validate:"required,gte=5"`
	Description    string              `json:"description" bson:"description" validate:"required,gte=5"`
	Author         *users.User         `json:"author" bson:"author"`
	ImageID        *primitive.ObjectID `json:"image" bson:"image"`
	CreatedAt      time.Time           `json:"createdAt" bson:"createdAt"`
	Views          int                 `json:"views,omitempty" bson:"views,omitempty"`
	MemberRequests []*users.User       `json:"membersRequest,omitempty" bson:"membersRequest,omitempty"`
	Members        []*users.User       `json:"members,omitempty" bson:"members,omitempty"`
}

func CreatePost(ctx context.Context, client *mongo.Client, authorID string, post *Post) (*Post, error) {
	coll := client.Database("codev").Collection("posts")

	author, err := users.GetUser(ctx, client, authorID)
	if err != nil {
		return nil, err
	}
	post.Author = author
	post.CreatedAt = time.Now()
	post.Views = 0
	post.Members = nil

	insertRes, err := coll.InsertOne(ctx, post)
	if err != nil {
		return nil, err
	}

	id, ok := insertRes.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("cannot assert id type")
	}
	post.ID = &id

	return post, nil
}

func GetPosts(ctx context.Context, client *mongo.Client, offset, limit int) ([]*Post, error) {
	coll := client.Database("codev").Collection("posts")
	var posts = make([]*Post, 0)
	cur, err := coll.Find(ctx, bson.D{}, options.Find().SetSkip(int64(offset)).SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var post Post
		err := cur.Decode(&post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func GetPost(ctx context.Context, client *mongo.Client, id string) (*Post, error) {
	coll := client.Database("codev").Collection("posts")
	var post Post
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrPostNotFound
	}
	singleRes := coll.FindOne(ctx, bson.D{
		{
			Key:   "_id",
			Value: objID,
		},
	})
	if singleRes.Err() != nil {
		if singleRes.Err() == mongo.ErrNoDocuments {
			return nil, ErrPostNotFound
		}
		return nil, singleRes.Err()
	}
	err = singleRes.Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	return &post, nil
}

func AddMember(ctx context.Context, client *mongo.Client, postID string, userID string) error {
	post, err := GetPost(ctx, client, postID)
	if err != nil {
		return err
	}

	userObj, err := users.GetUser(ctx, client, userID)
	if err != nil {
		return auth.ErrWrongToken
	}
	if userObj.ID.Hex() == post.Author.ID.Hex() {
		return ErrMemberIsAuthor
	}
	postsColl := client.Database("codev").Collection("posts")
	if post.Members == nil {
		post.Members = make([]*users.User, 0)
	} else {
		for _, member := range post.Members {
			if member.ID.Hex() == userObj.ID.Hex() {
				return ErrMemberAlreadyExists
			}
		}
	}
	post.Members = append(post.Members, userObj)
	_, err = postsColl.ReplaceOne(ctx, bson.D{
		{
			Key:   "_id",
			Value: post.ID,
		},
	}, post)
	if err != nil {
		return err
	}

	return nil
}

func DeleteMember(ctx context.Context, client *mongo.Client, postID string, userID string) error {
	post, err := GetPost(ctx, client, postID)
	if err != nil {
		return err
	}

	userObj, err := users.GetUser(ctx, client, userID)
	if err != nil {
		return auth.ErrWrongToken
	}

	postsColl := client.Database("codev").Collection("posts")
	if post.Members == nil {
		return ErrMebmerNotExists
	}
	deleted := false
	for i, member := range post.Members {
		if member.ID.Hex() == userObj.ID.Hex() {
			post.Members[i] = post.Members[len(post.Members)-1]
			post.Members = post.Members[:len(post.Members)-1]
			deleted = true
			break
		}
	}
	if !deleted {
		return ErrMebmerNotExists
	}
	_, err = postsColl.ReplaceOne(ctx, bson.D{
		{
			Key:   "_id",
			Value: post.ID,
		},
	}, post)
	if err != nil {
		return err
	}

	return nil
}

//func UploadImage(ctx context.Context, )
