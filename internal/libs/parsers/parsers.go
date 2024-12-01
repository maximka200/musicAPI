package parsers

import "strings"

func ParseInCouplets(text string) []string {
	trimmedText := strings.TrimSpace(text)

	couplets := strings.Split(trimmedText, "\n\n")

	return couplets
}

func JoinCouplets(couplets []string) string {
	return strings.Join(couplets, "\n\n")
}
