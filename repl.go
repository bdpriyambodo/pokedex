package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Pokedex > ")

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			break
		}

		fmt.Println("Your command was:", cleanInput(line)[0])

		fmt.Print("Pokedex > ")
	}
}

func cleanInput(text string) []string {
	// separator := " "
	// result := strings.Split(text, separator)
	lowerCaseText := strings.ToLower((text))
	result := strings.Fields(lowerCaseText)
	// fmt.Println("Result: ", result)
	// fmt.Println("Length: ", len(result))
	return result
}
