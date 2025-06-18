package domain

import "fmt"

func BatchUpdateMsg(username string, state string) string {
	return fmt.Sprintf("%v %v your request.", username, state)
}

func BatchUpdateReasonMsg(username string, state string, reason string) string {
	return fmt.Sprintf("%v %v your request: %v.", username, state, reason)
}
