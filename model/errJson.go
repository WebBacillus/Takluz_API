package model

type errorResponse struct {
	Error string `json:"error"` // Use "error" instead of "message"
}
