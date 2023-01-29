package utils

import "strings"

func CheckAndFormatURL(url string) string {
	if strings.HasPrefix(url, "https://") {
		return strings.Replace(url, "https://", "wss://", 1)
	} else if strings.HasPrefix(url, "http://") {
		return strings.Replace(url, "http://", "ws://", 1)
	}
	return ""
}
