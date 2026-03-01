package main

import (
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	// separator := " "
	// result := strings.Split(text, separator)
	lowerCaseText := strings.ToLower((text))
	result := strings.Fields(lowerCaseText)
	fmt.Println("Result: ", result)
	fmt.Println("Length: ", len(result))
	return result
}
