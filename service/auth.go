package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/Hunnnn77/hello/model"
	"github.com/Hunnnn77/hello/response"
	"github.com/Hunnnn77/hello/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Colls struct {
	UserCollection    *mongo.Collection
	ActressCollection *mongo.Collection
}

func (c Colls) Signup(w http.ResponseWriter, r *http.Request, body *model.SignUp) {
	if ok := util.ExistingCookie(r); ok {
		util.DelCookie(w)
	}
	filter := c.GetFilterByEmail(body.Email)
	if u := c.IsExisting(filter); u != nil {
		e := model.HttpError{
			Code:    http.StatusBadRequest,
			Message: util.Capitalize("signup: existing user"),
		}
		response.ThrowError(w, e)
		return
	}

	signupBody := model.SignUp{
		Email:     body.Email,
		Password:  body.Password,
		CreatedAt: util.NowInTimezone,
		Rt:        nil,
	}

	if _, err := c.UserCollection.InsertOne(context.TODO(), signupBody); err != nil {
		e := model.HttpError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
		response.ThrowError(w, e)
		return
	}

	response.ThrowOk[bool](w, model.HttpOk[bool, string]{
		Ok: true,
	})
}

func (c Colls) Login(w http.ResponseWriter, r *http.Request, body *model.LogIn) {
	filter := c.GetFilterByEmail(body.Email)

	if user := c.IsExisting(filter); user == nil {
		response.ThrowError(w, model.HttpError{
			Code:    http.StatusNotFound,
			Message: util.Capitalize("login: no user"),
		})
		return
	} else if user.Password != body.Password {
		response.ThrowError(w, model.HttpError{
			Code:    http.StatusNotFound,
			Message: util.Capitalize("login: invalid password"),
		})
		return
	}

	atExp, rtExp, cookieExp := util.GetExpiration()
	atToken, rtToken := util.GenToken(body.Email, atExp), util.GenToken(body.Email, rtExp)

	if err := c.UpdateRt(body.Email, &rtToken); err != nil {
		response.ThrowError(w, model.HttpError{
			Code:    http.StatusNotFound,
			Message: util.Capitalize("login: update error"),
		})
		return
	}

	if util.PlatForm == util.WEB {
		util.ThrowCookie(w, atToken, cookieExp)
	} else {
		w.Header().Set("jwt", atToken)
	}
	response.ThrowOk[bool](w, model.HttpOk[bool, string]{
		Ok: true,
	})
}

func (c Colls) Logout(w http.ResponseWriter, token string) {
	if email, err := util.GetEmail(util.EMAIL, token); err != nil {
		response.ThrowError(w, model.HttpError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	} else {
		if err := c.UpdateRt(*email, nil); err != nil {
			response.ThrowError(w, model.HttpError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}
		if util.PlatForm == util.WEB {
			util.DelCookie(w)
		}
		response.ThrowOk[bool](w, model.HttpOk[bool, string]{
			Ok: true,
		})
	}
}

func (c Colls) GetFilterByEmail(email string) bson.D {
	filter := bson.D{
		{Key: util.EMAIL, Value: email},
	}
	return filter
}

func (c Colls) IsExisting(filter bson.D) *model.SignUp {
	user := new(model.SignUp)
	if err := c.UserCollection.FindOne(context.TODO(), filter).Decode(user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
	} else {
		return user
	}
	return nil
}

func (c Colls) UpdateRt(email string, rtToken *string) error {
	filter := c.GetFilterByEmail(email)
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "rt", Value: rtToken},
			},
		},
	}
	_, err := c.UserCollection.UpdateOne(context.TODO(), filter, update)
	return err
}
