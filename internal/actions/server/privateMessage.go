package server

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func extractReceivers(message string) []string {
	header := strings.Split(message, "\n")[0]
	regex := regexp.MustCompile(`Private message, length=\d+ to ([\w,]+)`)
	usernames := regex.FindStringSubmatch(header)
	if len(usernames) < 2 {
		return nil
	}
	return strings.Split(usernames[1], ",")
}

func SendPrivateMessage(username string, message string, actions Actions) error {
	body, err := extractBody(message)
	if err != nil {
		return err
	}

	receivers := extractReceivers(message)

	if len(receivers) == 0 {
		return errors.New("no receivers found")
	}

	msg := fmt.Sprintf("Private message, length=%d from %s to %s:\n\r%s",
		len(body),
		username,
		strings.Join(receivers, ","),
		body)

	for _, receiverUsername := range receivers {
		actions.BroadcastTo(receiverUsername, msg)
	}
	return nil
}
