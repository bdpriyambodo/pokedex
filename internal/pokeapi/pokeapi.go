package pokeapi

import (
	"encoding/json"
	"fmt"
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

func PokeApiData(offset int, url string) {

	fullUrl := url + fmt.Sprintf("?offset=%d&limit=20", offset)

	res, err := http.Get(fullUrl)
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

	// fmt.Println("body:", string(body))
	// fmt.Println("err:", err)
	// fmt.Println("===============================")

	result := LocationAreaResponse{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("result:", result)
	// fmt.Println("===============================")

	for _, area := range result.Results {
		fmt.Println(area.Name)
	}

}
