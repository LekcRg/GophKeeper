package models

type ResponseError struct {
	Error string `json:"error"`
}

type Response struct {
	Message string `json:"message"`
}
