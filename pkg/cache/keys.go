package cache

import "fmt"

func SessionKey(sessionID string) string {
	return fmt.Sprintf("session:%s", sessionID)
}
