package request

import "time"

type AuthenticatedUserRequest struct {
	Username string
	Message  string
	Time     time.Time
}

func New(username string, message string) *AuthenticatedUserRequest {
	return &AuthenticatedUserRequest{Username: username, Message: message, Time: time.Now()}
}
