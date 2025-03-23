package server

import (
	"errors"
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
	regex := regexp.MustCompile(`length=(\d+)`)
	matches := regex.FindStringSubmatch(lines[0])
	length, err := strconv.Atoi(matches[1])
	if err != nil {
		return "", errors.New("failed to extract length")
	}
	if length > 2000 {
		return "", errors.New("message too long")
	}
	return strings.Join(lines[1:], "\n")[1 : length+1], nil
}
