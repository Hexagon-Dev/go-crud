package common

type HttpError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}
