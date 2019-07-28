package errors

import "github.com/misgorod/co-dev/common/errors"

var ErrNoPostID = errors.New("post's ID not specified", 400)
var ErrNoMemberID = errors.New("member's ID not specidied", 404)
var ErrPostNotFound = errors.New("post not found", 404)
var ErrMebmerNotExists = errors.New("member not exists", 404)
var ErrMemberIsAuthor = errors.New("member is not allowed to be an author of post", 400)
var ErrMemberAlreadyExists = errors.New("member already exists", 400)
var ErrNotAnAuthor = errors.New("user is not an author of post", 403)
var ErrNoFile = errors.New("file with this id not found", 404)
