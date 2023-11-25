package controller

import (
	"net/http"

	"github.com/Hunnnn77/hello/db"
	"github.com/Hunnnn77/hello/model"
	"github.com/Hunnnn77/hello/response"
	"github.com/Hunnnn77/hello/util"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if body, err := util.GenerateBody[*model.LogIn](r, new(model.LogIn)); err != nil {
		response.ThrowError(w, model.HttpError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	} else {
		db.Collections.Login(w, r, body)
	}
}

func HandleSignup(w http.ResponseWriter, r *http.Request) {
	if body, err := util.GenerateBody[*model.SignUp](r, new(model.SignUp)); err != nil {
		response.ThrowError(w, model.HttpError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	} else {
		db.Collections.Signup(w, r, body)
	}
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	if _, ctx := util.HandleContext(r, nil); ctx == nil {
		response.ThrowError(w, model.HttpError{
			Code:    http.StatusNotFound,
			Message: util.Capitalize("logout: not found context")},
		)
	} else {
		if token, ok := ctx.(string); ok {
			db.Collections.Logout(w, token)
		}
	}
}
