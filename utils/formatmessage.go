package utils

import (
	"fmt"
	"time"
)

// formatMessage formats a chat message by adding the current timestamp and client name.
func formatMessage(name, message string) string {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s][%s]:%s", currentTime, name, message)
}
