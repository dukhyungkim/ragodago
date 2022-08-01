package util

import "strings"

func ParseName(resourceURL string) string {
	s1 := strings.Split(resourceURL, ":")
	s2 := strings.Split(s1[0], "/")
	return s2[len(s2)-1]
}
