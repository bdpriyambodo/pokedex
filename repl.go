package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bdpriyambodo/pokedexcli/internal/pokeapi"
	"github.com/bdpriyambodo/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(existingcache *pokecache.Cache, offset int, url string) error
}

func startRepl() {

	existingcache := pokecache.NewCache(5 * time.Second)

	url := "https://pokeapi.co/api/v2/location-area/"
	offset := 0
	next_offset := 0

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
		"map": {
			name:        "map",
			description: "Show 20 area locations from the PokeAPI",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Show previous 20 area locations from the PokeAPI",
			callback:    commandMap,
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
			if !(cmd.name == "map") && !(cmd.name == "mapb") {
				commandList[cmd.name].callback(existingcache, offset, url)
			} else if cmd.name == "map" {
				offset = next_offset
				commandList[cmd.name].callback(existingcache, offset, url)
				next_offset += 20
			} else if cmd.name == "mapb" {
				if offset <= 0 {
					fmt.Println("You're already on the first page")
				} else {
					offset -= 20
					commandList[cmd.name].callback(existingcache, offset, url)
					next_offset -= 20
				}
			}
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

func commandExit(existingcache *pokecache.Cache, offset int, url string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(existingcache *pokecache.Cache, offset int, url string) error {
	fmt.Println("Please type available commands to use the Pokedex")
	return nil
}

func commandMap(existingcache *pokecache.Cache, offset int, url string) error {
	// fmt.Println("Here are the available area locations:")

	if offset >= 0 {
		fullUrl := url + fmt.Sprintf("?offset=%d&limit=20", offset)
		val, ok := existingcache.Get(fullUrl)

		if ok {
			names, _ := pokeapi.PokeApiData(val)
			for _, name := range names {
				fmt.Println(name)
			}
		} else {
			client := &http.Client{}
			body, _ := pokeapi.PokeApiRaw(client, offset, url)
			names, _ := pokeapi.PokeApiData(body)
			for _, name := range names {
				fmt.Println(name)
			}
			existingcache.Add(fullUrl, body)
		}

	} else {
		fmt.Println("You're already on the first page")
	}
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
