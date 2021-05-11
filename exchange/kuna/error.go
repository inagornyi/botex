package kuna

import "fmt"

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type KunaError struct {
	Error *Error `json:"error"`
}

func (e Error) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}
