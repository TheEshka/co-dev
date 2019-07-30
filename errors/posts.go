package errors

var ErrNoPostID = New("post's ID not specified", 400)
var ErrNoMemberID = New("member's ID not specidied", 404)
var ErrPostNotFound = New("post not found", 404)
var ErrMebmerNotExists = New("member not exists", 404)
var ErrMemberIsAuthor = New("member is not allowed to be an author of post", 400)
var ErrMemberAlreadyExists = New("member already exists", 400)
var ErrNotAnAuthor = New("user is not an author of post", 403)

