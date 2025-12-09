package cache

import "fmt"

func SessionKey(sessionID string) string {
	return fmt.Sprintf("session:%s", sessionID)
}

func UserKey(userID string) string {
	return fmt.Sprintf("user:%s", userID)
}

func TokenKey(token string) string {
	return fmt.Sprintf("token:%s", token)
}
