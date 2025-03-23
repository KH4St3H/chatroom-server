package request

import (
	"github.com/kh4st3h/chatroom-server/internal/constants"
	"strings"
	"time"
)

type AuthenticatedUserRequest struct {
	Username string
	Message  string
	Type     int
	Time     time.Time
}

func ClassifyType(message string) int {
	if message == constants.FETCH_ATTENDEES_MSG {
		return constants.FETCH_ATTENDEES_TYPE
	}
	if strings.HasPrefix(message, "Public message") {
		return constants.PUBLIC_MESSAGE_TYPE
	}
	if strings.HasPrefix(message, "Private message") {
	}
	if strings.HasPrefix(message, "Bye.") {
		return constants.BYE_MESSAGE_TYPE
	}
	return constants.UNKOWN_TYPE
}

func New(username string, message string) *AuthenticatedUserRequest {
	return &AuthenticatedUserRequest{Username: username, Message: message, Time: time.Now(),
		Type: ClassifyType(message)}
}
