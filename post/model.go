package post

import (
	"context"
	"io"
	"time"

	"github.com/misgorod/co-dev/common/errors"
	perrors "github.com/misgorod/co-dev/post/errors"
	"github.com/misgorod/co-dev/users"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
	ID             primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Title          string              `json:"title" bson:"title" validate:"required,gte=5"`
	Subject        string              `json:"subject" bson:"subject" validate:"required,gte=5"`
	Description    string              `json:"description" bson:"description" validate:"required,gte=5"`
	Author         *users.User         `json:"author" bson:"author"`
	ImageID        *primitive.ObjectID `json:"image,omitempty" bson:"image,omitempty"`
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
		return nil, errors.ErrAssertID
	}
	post.ID = id
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
		return nil, perrors.ErrPostNotFound
	}
	singleRes := coll.FindOne(ctx, bson.D{
		{
			Key:   "_id",
			Value: objID,
		},
	})
	err = singleRes.Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, perrors.ErrPostNotFound
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
		return errors.ErrWrongToken
	}
	if userObj.ID.Hex() == post.Author.ID.Hex() {
		return perrors.ErrMemberIsAuthor
	}
	postsColl := client.Database("codev").Collection("posts")
	if post.MemberRequests == nil {
		post.MemberRequests = make([]*users.User, 0)
	} else {
		for _, member := range post.MemberRequests {
			if member.ID.Hex() == userObj.ID.Hex() {
				return perrors.ErrMemberAlreadyExists
			}
		}
		for _, member := range post.Members {
			if member.ID.Hex() == userObj.ID.Hex() {
				return perrors.ErrMemberAlreadyExists
			}
		}
	}
	post.MemberRequests = append(post.MemberRequests, userObj)
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

func ApproveMember(ctx context.Context, client *mongo.Client, postID string, authorID string, memberID string) error {
	coll := client.Database("codev").Collection("posts")
	post, err := GetPost(ctx, client, postID)
	if err != nil {
		return err
	}
	if post.Author.ID.Hex() != authorID {
		return perrors.ErrNotAnAuthor
	}
	for _, member := range post.Members {
		if member.ID.Hex() == memberID {
			return perrors.ErrMemberAlreadyExists
		}
	}
	deleted := false
	var memberObj *users.User
	for i, member := range post.MemberRequests {
		if member.ID.Hex() == memberID {
			memberObj = post.MemberRequests[i]
			post.MemberRequests[i] = post.MemberRequests[len(post.MemberRequests)-1]
			post.MemberRequests = post.MemberRequests[:len(post.MemberRequests)-1]
			deleted = true
			break
		}
	}
	if !deleted {
		return perrors.ErrMebmerNotExists
	}
	post.Members = append(post.Members, memberObj)
	_, err = coll.ReplaceOne(ctx, bson.D{
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

func DeleteMemberSelf(ctx context.Context, client *mongo.Client, postID string, userID string) error {
	post, err := GetPost(ctx, client, postID)
	if err != nil {
		return err
	}
	postsColl := client.Database("codev").Collection("posts")
	if post.Members == nil {
		return perrors.ErrMebmerNotExists
	}
	deleted := false
	for i, member := range post.Members {
		if member.ID.Hex() == userID {
			post.Members[i] = post.Members[len(post.Members)-1]
			post.Members = post.Members[:len(post.Members)-1]
			deleted = true
			break
		}
	}
	if !deleted {
		return perrors.ErrMebmerNotExists
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

func DeleteMember(ctx context.Context, client *mongo.Client, postID string, authorID string, memberID string) error {
	post, err := GetPost(ctx, client, postID)
	if err != nil {
		return err
	}
	if post.Author.ID.Hex() != authorID {
		return perrors.ErrNotAnAuthor
	}
	postsColl := client.Database("codev").Collection("posts")
	if post.Members == nil {
		return perrors.ErrMebmerNotExists
	}
	deleted := false
	for i, member := range post.Members {
		if member.ID.Hex() == memberID {
			post.Members[i] = post.Members[len(post.Members)-1]
			post.Members = post.Members[:len(post.Members)-1]
			deleted = true
			break
		}
	}
	if !deleted {
		return perrors.ErrMebmerNotExists
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

func UploadImage(ctx context.Context, client *mongo.Client, reader io.Reader, post *Post) (*file, error) {
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
	filter := bson.D{{"_id", post.ID}}
	update := bson.D{{"$set", bson.D{
		{"image", fileID},
	}}}
	_, err = col.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	post.ImageID = &fileID

	return &file{fileID}, nil
}

func DownloadImage(ctx context.Context, client *mongo.Client, id string, writer io.Writer) error {
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
			return perrors.ErrNoFile
		}
		return err
	}
	return nil
}
