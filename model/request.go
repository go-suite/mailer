package model

type Request struct {
	Message        message         `json:"message"`
	Authentication *authentication `json:"authentication,omitempty"`
}
