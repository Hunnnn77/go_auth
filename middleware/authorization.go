package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Hunnnn77/hello/db"
	"github.com/Hunnnn77/hello/model"
	"github.com/Hunnnn77/hello/response"
	"github.com/Hunnnn77/hello/util"
	"github.com/go-chi/jwtauth/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch util.PlatForm {
		case util.WEB:
			cookie, err := util.ParseCookie(r)
			if err != nil {
				response.ThrowError(w, model.HttpError{
					Code:    http.StatusUnauthorized,
					Message: err.Error(),
				})
				return
			}
			if len(cookie.Value) == 0 {
				response.ThrowError(w, model.HttpError{
					Code:    http.StatusUnauthorized,
					Message: util.Capitalize("empty cookie"),
				})
				return
			}
			if err := refreshCookieIf(w, r, next, cookie.Value); err != nil {
				response.ThrowError(w, model.HttpError{
					Code:    http.StatusUnauthorized,
					Message: util.Capitalize(err.Error()),
				})
				return
			}
		case util.MOBILE:
			if token, err := readToken(w, r); err != nil {
				response.ThrowError(w, model.HttpError{
					Code:    http.StatusUnauthorized,
					Message: util.Capitalize(err.Error()),
				})
				return
			} else {
				if err := refreshTokenIf(w, r, next, *token); err != nil {
					response.ThrowError(w, model.HttpError{
						Code:    http.StatusUnauthorized,
						Message: util.Capitalize(err.Error()),
					})
					return
				}
			}
		default:
			panic("unreachable")
		}
	})
}

func refreshCookieIf(w http.ResponseWriter, r *http.Request, next http.Handler, token string) error {
	email, err := util.GetEmail(util.EMAIL, token)
	if err != nil {
		return err
	}

	if err = util.VerifyToken(token); err == nil {
		return nil
	}

	if !errors.Is(err, jwtauth.ErrExpired) {
		return err
	}
	c := db.Collections
	filter := c.GetFilterByEmail(*email)
	user := c.IsExisting(filter)

	switch {
	case user == nil:
		return mongo.ErrNoDocuments
	case user.Rt == nil:
		return util.ErrEmptyRtToken
	}

	if err = util.VerifyToken(*user.Rt); err != nil {
		if !errors.Is(err, jwtauth.ErrExpired) {
			return err
		}

		_, rtExp, _ := util.GetExpiration()
		rt := util.GenToken(*email, rtExp)
		if err := c.UpdateRt(*email, &rt); err != nil {
			return err
		}
	}
	refreshCookie(w, r, email)
	req, _ := util.HandleContext(r, token)
	next.ServeHTTP(w, req)
	return nil
}

func refreshCookie(w http.ResponseWriter, r *http.Request, email *string) {
	atExp, _, cookieExp := util.GetExpiration()
	at := util.GenToken(*email, atExp)
	util.DelCookie(w)
	util.ThrowCookie(w, at, cookieExp)
}

func refreshTokenIf(w http.ResponseWriter, r *http.Request, next http.Handler, at string) error {
	if err := util.VerifyToken(at); err == nil {
		req, _ := util.HandleContext(r, at)
		next.ServeHTTP(w, req)
		return nil
	} else {
		if !errors.Is(err, jwtauth.ErrExpired) {
			return jwtauth.ErrUnauthorized
		}
		if err := util.RemovePreviousToken(at); err != nil {
			return err
		}
		if email, err := util.GetEmail(util.EMAIL, at); err != nil {
			return err
		} else {
			c := db.Collections
			filter := c.GetFilterByEmail(*email)
			user := c.IsExisting(filter)
			if user == nil {
				return util.ErrUserNotExist
			}
			if err = util.VerifyToken(*user.Rt); err != nil {
				if !errors.Is(err, jwtauth.ErrExpired) {
					return jwtauth.ErrUnauthorized
				}

				_, rtExp, _ := util.GetExpiration()
				rt := util.GenToken(*email, rtExp)
				if err = c.UpdateRt(*email, &rt); err != nil {
					return util.ErrUpdatingRt
				}
			}
			atExp, _, _ := util.GetExpiration()
			newToken := util.GenToken(*email, atExp)

			refreshHeader(w, r, newToken)
			tokenInHeader, err := readToken(w, r)
			if err != nil {
				return err
			}

			if *tokenInHeader != "" {
				req, _ := util.HandleContext(r, *tokenInHeader)
				next.ServeHTTP(w, req)
				return nil
			} else {
				return util.ErrEmptyAtToken
			}
		}
	}
}

func readToken(w http.ResponseWriter, r *http.Request) (*string, error) {
	token := strings.Split(r.Header.Get("Authorization"), " ")
	if len(token) != 2 || token[1] == "" {
		return nil, util.ErrEmptyHeader
	}
	return &token[1], nil
}

func refreshHeader(w http.ResponseWriter, r *http.Request, at string) {
	key := "Authorization"
	r.Header.Del(key)
	r.Header.Set(key, "Bearer "+at)
	w.Header().Set("jwt", at)
}
