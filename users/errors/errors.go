package errors

import "github.com/misgorod/co-dev/common/errors"

var ErrNoID = errors.New("user's ID not specified", 400)
var ErrUserNotExists = errors.New("user with this id doesn't exist", 404)
