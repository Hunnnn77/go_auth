package model

type Authorized struct {
	Email string  `json:"email,omitempty"`;
	At    *string `json:"at,omitempty"`
}

type Ok interface {
	Authorized | bool | string
}

type HttpOk[T Ok, U Ok] struct {
	Ok T  `json:"ok,omitempty"`
	At *U `json:"at,omitempty"`
}

type HttpError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
