package main

import (
	"strconv"
	"unicode"
)

// To convert M, B, T in a number to a float64 (ie market cap)
func parseNum(s string) float64 {
	var number string
	var result float64
	for _, char := range s {
		if !unicode.IsLetter(char) {
			number += string(char)
		}
	}

	if len(number) == len(s) {
		res, _ := strconv.ParseFloat(number, 64)
		return res
	} else {
		res, _ := strconv.ParseFloat(number, 64)
		switch s[len(s)-1] {
		case 'M':
			result = res * 1000000
		case 'B':
			result = res * 1000000000
		case 'T':
			result = res * 1000000000
		}
	}
	return result
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
