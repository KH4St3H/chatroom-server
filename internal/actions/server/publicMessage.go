package server

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Actions interface {
	Broadcast(sender string, message string)
	BroadcastTo(username string, msg string)
}

func extractBody(msg string) (string, error) {
	lines := strings.Split(msg, "\n")
	regex := regexp.MustCompile(`Public message, length=(\d+)`)
	matches := regex.FindStringSubmatch(lines[0])
	length, err := strconv.Atoi(matches[1])
	if err != nil {
		return "", errors.New("failed to extract length")
	}
	if length > 2000 {
		return "", errors.New("message too long")
	}
	return strings.Join(lines[1:], "\n")[:length+1], nil
}

func SendPublicMessage(username string, message string, actions Actions) error {
	body, err := extractBody(message)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("Public message from %s, length=%d\n\r%s", username, len(body), body)
	actions.Broadcast(username, msg)
	return nil
}
