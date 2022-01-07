package entity

import "errors"

var ErrEmailExist = errors.New("email already exist")
var ErrUsernameExist = errors.New("username already exist")

var ErrUserNotFound = errors.New("user not found")
var ErrBroadcastNotFound = errors.New("broadcast not found")

var ErrPasswordIncorrect = errors.New("password incorrect")
var ErrPasswordRepeatIncorrect = errors.New("password is not same as repeated password")

var ErrAuthTokenInvalid = errors.New("auth token invalid")
var ErrAuthTokenExpired = errors.New("auth token expired")
