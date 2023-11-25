package util

import (
	"errors"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

type ContextKey string

const (
	DEV int = iota
	PROD
	WEB
	MOBILE
	JWT       ContextKey = "JWT"
	PORT                 = 3000
	EMAIL                = "email"
	EXP                  = "exp"
	COOKIEKEY            = "jwt"
	ATTOKEN              = "at"
	RTTOKEN              = "rt"
)

var (
	MODE                      = DEV
	PlatForm                  = MOBILE
	JwtAuth  *jwtauth.JWTAuth = nil
)

var (
	ErrInvalidKey     = errors.New("invalid key")
	ErrConversion     = errors.New("failed to conversion")
	ErrEmptyRtToken   = errors.New("logged out")
	ErrEmptyAtToken   = errors.New("empty at token")
	ErrEmptyHeader    = errors.New("header - authorization is empty")
	ErrInvalidContext = errors.New("invalid context")
	ErrUserNotExist   = errors.New("not existing user")
	ErrUpdatingRt     = errors.New("updating rt error")
)

func isDev(key int) (time.Time, time.Time) {
	switch key {
	case DEV:
		return time.Now().Add(time.Minute), time.Now().Add(time.Minute * 2)
	case PROD:
		return time.Now().Add(time.Minute * 2), time.Now().Add(time.Minute * 3)
	default:
		panic("unreachable")
	}
}

func GetExpiration() (at int64, rt int64, cookie time.Time) {
	a, r := isDev(MODE)
	at = a.Unix()
	rt = r.Unix()
	cookie = r
	return
}
