package post

import "errors"

var ErrPostNotFound = errors.New("post not found")
var ErrMebmerNotExists = errors.New("member not exists")
var ErrMemberIsAuthor = errors.New("member is not allowed to be an author of post")
var ErrMemberAlreadyExists = errors.New("member already exists")
