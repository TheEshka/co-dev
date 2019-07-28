package errors

import "github.com/misgorod/co-dev/common/errors"

var ErrUserExists = errors.New("user already exists", 400)
var ErrUserNotExists = errors.New("user with this id doesn't exist", 404)
var ErrWrongCreds = errors.New("wrong login or password", 401)