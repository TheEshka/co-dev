package users

import "errors"

var ErrUserExists = errors.New("user already exists")
var ErrUserNotExists = errors.New("user with this id doesn't exist")
