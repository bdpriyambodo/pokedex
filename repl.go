package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func startRepl() {

	commandList := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "How to use the Pokedex",
			callback:    commandHelp,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	for key := range commandList {
		fmt.Println(commandList[key].name, ":", commandList[key].description)
	}

	fmt.Print("\nPokedex > ")

	for scanner.Scan() {

		line := scanner.Text()

		word := cleanInput(line)[0]

		cmd, exist := commandList[word]

		if exist {
			commandList[cmd.name].callback()
		} else {
			fmt.Println("Unknown command")
		}

		// switch word {
		// case "exit":
		// 	commandList["exit"].callback()
		// case "help":
		// 	commandList["help"].callback()
		// default:
		// 	fmt.Println("Unknown command")
		// }

		fmt.Print("\nPokedex > ")
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Please type available commands to use the Pokedex")
	return nil
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
