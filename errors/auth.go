package errors

var ErrUserExists = New("user already exists", 400)
var ErrUserNotExists = New("user with this id doesn't exist", 404)
var ErrWrongCreds = New("wrong login or password", 401)
