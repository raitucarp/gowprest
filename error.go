package gowprest

import "fmt"

type WPRestError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Data    struct {
		Status int `json:"status"`
	} `json:"data"`
}

func (e *WPRestError) Error() string {
	return fmt.Sprintf("[%d][%s] %s", e.Data.Status, e.Code, e.Message)
}
