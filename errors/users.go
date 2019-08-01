package errors

var ErrNoID = New("user's ID not specified", 400)
var ErrUserNotFound = New("user not found", 404)
