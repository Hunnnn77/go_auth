package util

import (
	"net/http"
	"time"
)

func newCookie(key, value string, exp time.Time) http.Cookie {
	return http.Cookie{
		Name:     key,
		Value:    value,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}
}

func delCookie(key string) http.Cookie {
	return http.Cookie{
		Name:     key,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}
}

func ThrowCookie(w http.ResponseWriter, atToken string, cookieExp time.Time) {
	cookie := newCookie(COOKIEKEY, atToken, cookieExp)
	http.SetCookie(w, &cookie)
}

func DelCookie(w http.ResponseWriter) {
	cookie := delCookie(COOKIEKEY)
	http.SetCookie(w, &cookie)
}

func ExistingCookie(r *http.Request) bool {
	if _, err := r.Cookie(COOKIEKEY); err != nil {
		return false
	}
	return true
}

func ParseCookie(r *http.Request) (*http.Cookie, error) {
	if ok := ExistingCookie(r); !ok {
		return nil, http.ErrNoCookie
	}

	jwtCookie, err := r.Cookie(COOKIEKEY)
	if err != nil {
		return nil, err
	}
	return jwtCookie, nil
}
