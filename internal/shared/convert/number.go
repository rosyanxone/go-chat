package convert

import (
	"fmt"
	"strconv"
)

func UintToString(value uint) string {
	return strconv.Itoa(int(value))
}

func IntToString(value int) string {
	return strconv.Itoa(value)
}

func StringToInt(value string) int {
	result, err := strconv.Atoi(value)

	if err != nil {
		fmt.Println("Invalid number:", err)
		return 0
	}

	return result
}
