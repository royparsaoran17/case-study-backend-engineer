package strx

import "strings"

// FirstToUpper conversion word first
// character to Upper case
func FirstToUpper(str string) string {
	if str == "" {
		return ""
	}
	return strings.ToUpper(str[0:1]) + str[1:]
}
