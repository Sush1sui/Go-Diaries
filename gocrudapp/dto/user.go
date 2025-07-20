package dto

type UserRequest struct {
	Name    string `json:"name,omitempty"`
	Age     int    `json:"age,omitempty"`
	Country string `json:"country,omitempty"`
}

type UserResponse struct {
	// Data can be any type, but typically it will be a struct or a slice of structs
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}