package convert

import (
	"regexp"
	"strings"
)

func NormalizePhoneNumber(phone string) string {
	re := regexp.MustCompile(`\D`)
	phone = re.ReplaceAllString(phone, "")

	switch {
	case strings.HasPrefix(phone, "620"):
		phone = "62" + phone[3:]

	case strings.HasPrefix(phone, "0"):
		phone = "62" + phone[1:]

	case strings.HasPrefix(phone, "62"):
		// Already normalized

	default:
		phone = "62" + phone
	}

	return phone
}
