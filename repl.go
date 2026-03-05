package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
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
	callback    func(pokedex map[string]*pokeapi.Pokemon, existingcache *pokecache.Cache, offset int, url string, input string) error
}

// type Pokemon struct {
// 	name   string
// 	caught int
// 	height int
// 	weight int
// 	hp int
// 	attack int
// 	defence int
// 	specialAttack int
// 	specialDefence int
// 	speed int
// 	typeTrait []string
// }

func startRepl() {

	existingcache := pokecache.NewCache(5 * time.Second)

	pokedex := make(map[string]*pokeapi.Pokemon)

	// testurl := "https://pokeapi.co/api/v2/location-area/"
	// fmt.Println(testurl)
	url := "https://pokeapi.co/api/v2/"
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
		"explore": {
			name:        "explore",
			description: "List all the Pokemon located in the specified area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch pokemon!",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Provide info on caught Pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all pokemons you have in Pokedex",
			callback:    commandPokedex,
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

		words := cleanInput(line)

		cmd, exist := commandList[words[0]]

		input := ""
		if len(words) > 1 {
			input = words[1]
		}

		if exist {
			if !(cmd.name == "map") && !(cmd.name == "mapb") {
				commandList[cmd.name].callback(pokedex, existingcache, offset, url, input)
			} else if cmd.name == "map" {
				offset = next_offset
				commandList[cmd.name].callback(pokedex, existingcache, offset, url, input)
				next_offset += 20
			} else if cmd.name == "mapb" {
				if offset <= 0 {
					fmt.Println("You're already on the first page")
				} else {
					offset -= 20
					commandList[cmd.name].callback(pokedex, existingcache, offset, url, input)
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

func commandExit(pokedex map[string]*pokeapi.Pokemon, existingcache *pokecache.Cache, offset int, url string, input string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(pokedex map[string]*pokeapi.Pokemon, existingcache *pokecache.Cache, offset int, url string, input string) error {
	fmt.Println("Please type available commands to use the Pokedex")
	return nil
}

func commandMap(pokedex map[string]*pokeapi.Pokemon, existingcache *pokecache.Cache, offset int, url string, input string) error {
	// fmt.Println("Here are the available area locations:")
	fullUrl := url + "location-area/" + fmt.Sprintf("?offset=%d&limit=20", offset)

	if offset >= 0 {
		val, ok := existingcache.Get(fullUrl)

		if ok {
			names, _ := pokeapi.PokeApiLocationArea(val)
			for _, name := range names {
				fmt.Println(name)
			}
		} else {
			client := &http.Client{}

			body, _ := pokeapi.PokeApiRaw(client, fullUrl)
			names, _ := pokeapi.PokeApiLocationArea(body)
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

func commandExplore(pokedex map[string]*pokeapi.Pokemon, existingcache *pokecache.Cache, offset int, url string, input string) error {
	fullUrl := url + "location-area/" + input

	val, ok := existingcache.Get(fullUrl)

	if ok {
		names, _ := pokeapi.PokeApiPokemonInArea(val)
		fmt.Println("Exploring" + input)
		fmt.Println("Found Pokemon:")
		for _, name := range names {
			fmt.Println("-", name)
		}
	} else {
		client := &http.Client{}
		body, _ := pokeapi.PokeApiRaw(client, fullUrl)
		names, _ := pokeapi.PokeApiPokemonInArea(body)
		fmt.Println("Exploring " + input)
		fmt.Println("Found Pokemon:")
		for _, name := range names {
			fmt.Println("-", name)
		}
		existingcache.Add(fullUrl, body)

	}

	return nil
}

func commandCatch(pokedex map[string]*pokeapi.Pokemon, existingcache *pokecache.Cache, offset int, url string, input string) error {
	fullUrl := url + "pokemon/" + input
	// fmt.Println(fullUrl)

	fmt.Printf("Throwing a Pokeball at %s...\n", input)

	randNumber := rand.IntN(250)
	var ptrPokemon *pokeapi.Pokemon

	val, ok := existingcache.Get(fullUrl)

	if ok {
		ptrPokemon, _ = pokeapi.PokeApiPokemonCatch(val)
	} else {
		client := &http.Client{}
		// body, _ := pokeapi.PokeApiRaw(client, fullUrl)

		body, err := pokeapi.PokeApiRaw(client, fullUrl)
		if err != nil {
			fmt.Println("There is no pokemon:", input)
			fmt.Println(err)

			return nil
		}

		ptrPokemon, _ = pokeapi.PokeApiPokemonCatch(body)
		existingcache.Add(fullUrl, body)

	}

	// fmt.Println("random number:", randNumber)
	// fmt.Println("base experience: ", ptrPokemon.BaseExperience)

	if randNumber > ptrPokemon.BaseExperience {
		fmt.Println(input, "was caught!")
		fmt.Println("You mau now inspect it with inspect command")
		pokemon, ok := pokedex[input]
		if ok {
			pokemon.Caught += 1

		} else {
			// newPokemon := pokeapi.Pokemon{
			// 	Name:   input,
			// 	Caught: 1,
			// }
			// pokedex[input] = &newPokemon
			pokedex[input] = ptrPokemon
			ptrPokemon.Caught += 1
		}
		// fmt.Println(input, pokedex[input].caught)
	} else {
		fmt.Println(input, "escaped!")
	}

	return nil

}

func commandInspect(pokedex map[string]*pokeapi.Pokemon, existingcache *pokecache.Cache, offset int, url string, input string) error {
	pokemon, ok := pokedex[input]
	if ok {
		fmt.Println("Name:", input)
		fmt.Printf("Height: %d\n", pokemon.Height)
		fmt.Printf("Weight: %d\n", pokemon.Weight)
		fmt.Println("Stats:")
		fmt.Printf("	-hp: %d\n", pokemon.Hp)
		fmt.Printf("	-attack: %d\n", pokemon.Attack)
		fmt.Printf("	-defence: %d\n", pokemon.Hp)
		fmt.Printf("	-special-attack: %d\n", pokemon.SpecialAttack)
		fmt.Printf("	-special-defence: %d\n", pokemon.SpecialDefence)
		fmt.Println(("Types:"))
		for _, text := range pokemon.TypeTrait {
			fmt.Printf("	-%s\n", text)
		}
		fmt.Printf("Caught: %d\n", pokemon.Caught)
	} else {
		fmt.Println("you have not caught that pokemon")
	}

	return nil
}

func commandPokedex(pokedex map[string]*pokeapi.Pokemon, existingcache *pokecache.Cache, offset int, url string, input string) error {
	fmt.Println("Your Pokedex:")
	for name := range pokedex {
		fmt.Printf(" -%s\n", name)
	}
	fmt.Println("")
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
