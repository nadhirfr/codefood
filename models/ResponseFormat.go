package models

type ResponseResult struct {
	Success bool        `json:"success" swaggertype:"bool"`
	Message string      `json:"message" swaggertype:"string"`
	Data    interface{} `json:"data" swaggertype:"object"`
}

type ResponseError struct {
	Success bool   `json:"success" swaggertype:"bool"`
	Message string `json:"message" swaggertype:"string"`
}
