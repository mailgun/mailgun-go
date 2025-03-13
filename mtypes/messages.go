package mtypes

// TODO(v5): return from Send()
type SendMessageResponse struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}
