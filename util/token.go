package util

import (
	"log"
	"reflect"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func RemovePreviousToken(token string) error {
	t, _ := jwtauth.VerifyToken(JwtAuth, token)
	exp := t.Expiration().Add(time.Minute).UnixMilli()
	if exp < time.Now().UnixMilli() {
		return jwtauth.ErrUnauthorized
	} else {
		return nil
	}
}

func GetEmail(key, token string) (*string, error) {
	t, _ := jwtauth.VerifyToken(JwtAuth, token)
	if val, ok := t.Get(key); !ok {
		return nil, jwt.ErrInvalidJWT()
	} else {
		switch v := val.(type) {
		case string:
			return &v, nil
		default:
			log.Printf("Failed Conversion |> value: %v | type: %v", v, reflect.TypeOf(v))
			return nil, ErrConversion
		}
	}
}

func VerifyToken(token string) error {
	_, err := jwtauth.VerifyToken(JwtAuth, token)
	if err != nil {
		return err
	}
	return nil
}

func GenToken(email string, exp int64) string {
	_, tokenString, _ := JwtAuth.Encode(map[string]any{
		"iat":   "domain",
		"email": email,
		"exp":   exp,
	})
	return tokenString
}
