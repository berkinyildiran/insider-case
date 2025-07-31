package message

type GenericSuccessResponse struct {
	Message string `json:"message"`
}

type GenericFailureResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
