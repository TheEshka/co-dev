package auth

import "errors"

var ErrUserExists = errors.New("user already exists")
var ErrWrongToken = errors.New("wrong token")
