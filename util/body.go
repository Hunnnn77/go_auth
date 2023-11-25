package util

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Hunnnn77/hello/model"
)

type RequestBody interface {
	*model.LogIn | *model.SignUp
}

func GenerateBody[T RequestBody](r *http.Request, m T) (T, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(r.Body)

	err = json.Unmarshal(body, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
