package models

type ResponseError struct {
	Error string `json:"error"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}

type Response struct {
	Message string `json:"message"`
}
