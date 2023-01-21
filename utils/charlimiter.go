package utils

import (
	"errors"
	"fmt"
	"strings"
)

// Limit characters of a string
func CheckCharLimit(msg string, limit int, codeblock bool) (string, error) {
	if msg == "" {
		return "", errors.New("empty msg is not allowed in CheckCharLimit")
	}
	if len(msg) >= limit {
		// Send full output to PrivateBin
		resp, err := PrivateBinPaste(msg)
		if err != nil {
			fmt.Printf("cannot create the paste: %s", err)
		}

		// Limit message output to not exeeed mattermost message limit - head
		//msg = msg[:strings.LastIndex(msg[:limit], " ")]
		// Limit message output to not exeeed mattermost message limit - tail
		msg = msg[len(msg)-limit:]
		msg += "\n"

		var result string
		if len(resp) > 0 {
			result = "Message limit exceeded. Click [here](" + resp + ") for full output.\n"
		} else {
			result = "Message limit exceeded.\n"
		}

		result += fmt.Sprintf("``` \n %s ```", msg)

		msg = strings.Replace(result, "\n\n", "\n", -1)
		return msg, nil
	}
	if codeblock {
		// Surround msg with graves (mattermost code format)
		result := fmt.Sprintf("``` \n %s ```", msg)
		msg = strings.Replace(result, "\n\n", "\n", -1)
	}

	return msg, nil
}
