package auth

import "errors"

var ErrWrongToken = errors.New("wrong token")
var errUserExists = errors.New("user already exists")
var errUserNotExists = errors.New("user with this id doesn't exist")
