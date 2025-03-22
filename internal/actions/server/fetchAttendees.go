package server

import (
	"strings"
)

func FetchAttendees(username string, userList []string, action Actions) {
	msg := "Here are the list of attendees:\n" + strings.Join(userList, ", ")
	action.BroadcastTo(username, msg)
}
