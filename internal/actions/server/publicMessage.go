package server

import (
	"fmt"
)

func SendPublicMessage(username string, message string, actions Actions) error {
	body, err := extractBody(message)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("Public message from %s, length=%d\n\r%s", username, len(body), body)
	actions.Broadcast(username, msg)
	return nil
}
