package convert

import "strconv"

func UintToString(value uint) string {
	return strconv.Itoa(int(value))
}

func IntToString(value int) string {
	return strconv.Itoa(value)
}
