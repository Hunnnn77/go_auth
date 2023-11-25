package response

import (
	"encoding/json"
	"net/http"

	"github.com/Hunnnn77/hello/model"
)

func ThrowOk[T model.Ok | bool, U string](w http.ResponseWriter, ok model.HttpOk[T, U]) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(ok)
}

func ThrowError(w http.ResponseWriter, err model.HttpError) {
	w.WriteHeader(err.Code)
	_ = json.NewEncoder(w).Encode(err)
}
