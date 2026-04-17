package service

type APIError struct {
	Status   int    `json:"-"`
	Message  string `json:"message"`
}

func (e *APIError) Error() string {
	return e.Message
}
