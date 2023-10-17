package entity

import "strings"

type BearerToken string

func (b BearerToken) GetToken() string {
	s := strings.Split(b.String(), "Bearer ")
	if len(s) > 1 {
		return s[1]
	}

	return ""
}

func (b BearerToken) String() string {
	return string(b)
}

func (b BearerToken) TokenEmpty() bool {
	return b.GetToken() == ""
}
