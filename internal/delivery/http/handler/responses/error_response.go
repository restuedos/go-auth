package responses

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code,omitempty"`
}
