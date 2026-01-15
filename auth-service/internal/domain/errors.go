package domain

import (
	"errors"
)

var ErrUserIdRequired = errors.New("user id is required")
var ErrUserNotFound = errors.New("user not found")

var ErrUserNameIsToShort = errors.New("user name is too short")
var ErrUserNameIsToLong = errors.New("user name is too long")

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrInvalidEmail = errors.New("invalid email")

var ErrPwdIsToShort = errors.New("password is too short")
var ErrPwdIsToLong = errors.New("password is too long")
var ErrPwdHashGeneration = errors.New("password hash generation")

var ErrInvalidCredentials = errors.New("invalid credentials")
