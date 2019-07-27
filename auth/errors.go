package auth

import "errors"

var ErrWrongToken = errors.New("wrong token")
var ErrUserExists = errors.New("user already exists")
var ErrUserNotExists = errors.New("user with this id doesn't exist")
var ErrAssertID = errors.New("cannot assert id")
