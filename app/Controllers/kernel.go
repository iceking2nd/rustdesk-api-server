package Controllers

type ResponseError struct {
	Error string `json:"error"`
}

type Response struct {
	Data string `json:"data"`
}
