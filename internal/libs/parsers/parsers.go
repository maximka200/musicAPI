package parsers

import (
	"fmt"
	"strings"
)

func ParseInCouplets(text string) []string {
	trimmedText := strings.TrimSpace(text)

	couplets := strings.Split(trimmedText, "\\n\\n")

	return couplets
}

func JoinCouplets(couplets []string) string {
	return strings.Join(couplets, "\n\n")
}

func StringDateForPsql(date string) (string, error) {
	arr := strings.Split(date, ".")
	if len(arr) != 3 {
		return "", fmt.Errorf("invalid date format: %s", date)
	}

	return fmt.Sprintf("%s-%s-%s", arr[2], arr[1], arr[0]), nil
}

func IsValidDate(date string) bool {
	arr := strings.Split(date, ".")
	if len(arr) != 3 {
		return false
	}
	if len(arr[0]) != 2 || len(arr[1]) != 2 || len(arr[2]) != 4 {
		return false
	}
	return true
}
