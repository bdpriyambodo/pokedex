package pokeapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type LocationAreaResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type PokemonInAreaResponse struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func PokeApiRaw(client *http.Client, fullUrl string) ([]byte, error) {

	// fullUrl := url + fmt.Sprintf("?offset=%d&limit=20", offset)

	res, err := client.Get(fullUrl)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("res:", res)
	// fmt.Println("err:", err)
	// fmt.Println("===============================")

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	return body, err

	// fmt.Println("body:", string(body))
	// fmt.Println("err:", err)
	// fmt.Println("===============================")

}

func PokeApiLocationArea(body []byte) ([]string, error) {
	result := LocationAreaResponse{}

	err := json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("result:", result)
	// fmt.Println("===============================")

	// for _, area := range result.Results {
	// 	fmt.Println(area.Name)
	// }

	// var names []string
	names := []string{}
	for _, area := range result.Results {
		names = append(names, area.Name)
	}

	return names, nil
}

func PokeApiPokemonInArea(body []byte) ([]string, error) {
	result := PokemonInAreaResponse{}

	err := json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("result:", result)
	// fmt.Println("===============================")

	// for _, area := range result.Results {
	// 	fmt.Println(area.Name)
	// }

	// var names []string
	names := []string{}
	for _, encounter := range result.PokemonEncounters {
		names = append(names, encounter.Pokemon.Name)
	}

	return names, nil
}
