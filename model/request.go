package model

import "time"

type LogIn struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type SignUp struct {
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Rt        *string   `json:"rt,omitempty"`
}
