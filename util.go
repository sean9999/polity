package polity

import (
	"fmt"
	"strings"
	"unicode"
)

func ifErr(err error, msg string) error {
	if err != nil {
		err = fmt.Errorf("%w: %s", err, msg)
	}
	return err
}

func isHex(str string) bool {
	// Check if the string starts with '0x' or '0X'
	if len(str) >= 2 && (str[:2] == "0x" || str[:2] == "0X") {
		str = str[2:]
	}

	// Check if the remaining string contains only valid hex digits
	if len(str) == 0 {
		return false
	}

	for _, char := range str {
		if !unicode.IsDigit(char) && !strings.ContainsRune("aAbBcCdDeEfF", char) {
			return false
		}
	}

	return true
}

func stringIsPubkey(str string) bool {
	//	ex: 1692c1beaa4021f60487acfaaf98fb5069c0b9bd9531585895009ca5730b5f62b4208154d0f5b39b37d62c7cc9ebff1e7f4d6425d68a154af5f6c93ab6e1a12e

	if len(str) != 128 {
		return false
	}

	return isHex(str)

}

func stringIsNickname(str string) bool {
	//	ex: hidden-butterfly

	//	must have a dash
	if !strings.Contains(str, "-") {
		return false
	}

	//	must be lowercase
	if str != strings.ToLower(str) {
		return false
	}

	//	must be greater than 4 chars
	if len(str) <= 4 {
		return false
	}

	//	must be less than 127 chars
	if len(str) > 127 {
		return false
	}

	return true
}
