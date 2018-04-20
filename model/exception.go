package model

// Exception holds the error message
type Exception struct {
	Error string `json:"error"`
}

// NewException takes error object and return an Exception type
func NewException(error error) Exception {
	return Exception{
		Error: error.Error(),
	}
}
