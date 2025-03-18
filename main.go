package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/CookieBorn/pokedexcli/internal/httphandels"
	"github.com/CookieBorn/pokedexcli/internal/pokecache"
)

var cache *pokecache.Cache
var clean []string
var pokedex map[string]statsPokemon

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cache = pokecache.NewCache(5 * time.Second)
	pokedex = make(map[string]statsPokemon)
	for true {
		fmt.Print("Pokedex >")
		scanner.Scan()
		command := scanner.Text()
		clean = cleanInput(command)
		calledCommand := false
		for com, cli := range callbackMap {
			if com == clean[0] {
				calledCommand = true
				cli.callback()
			}
		}
		if calledCommand == false {
			fmt.Print("Unknown command\n")
		}

	}
}

var mapNext string
var mapPrevious any

func cleanInput(text string) []string {
	textStrip := strings.TrimSpace(text)
	textStrip = strings.ToLower(textStrip)
	result := strings.Split(textStrip, " ")
	return result
}

func commandExit() error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, cli := range callbackMap {
		fmt.Printf("%s: %s\n", cli.name, cli.description)
	}
	return nil
}

func commandMap() error {
	url := ""
	var locations Location
	if mapNext != "" {
		url = mapNext
	} else {
		url = "https://pokeapi.co/api/v2/location-area/"
	}
	body, ok := cache.Get(url)
	if ok {
		if err := json.Unmarshal(body, &locations); err != nil {
			fmt.Printf("locations error: %v", err)
			return fmt.Errorf("Unmarshal error: %v", err)
		}
	} else {
		body, _ := httphandels.HTTPGet(url)
		cache.Add(url, body)

		if err := json.Unmarshal(body, &locations); err != nil {
			fmt.Printf("locations error: %v", err)
			return fmt.Errorf("Unmarshal error: %v", err)
		}
	}
	mapNext = locations.Next
	mapPrevious = locations.Previous
	for i := 0; i < len(locations.Results); i++ {
		fmt.Printf("%s\n", locations.Results[i].Name)
	}

	return nil
}

func commandMapB() error {
	url := ""
	var locations Location
	if mapPrevious != nil {
		url = mapPrevious.(string)
	} else {
		fmt.Printf("end of map\n")
		return nil
	}
	body, ok := cache.Get(url)
	if ok {
		if err := json.Unmarshal(body, &locations); err != nil {
			fmt.Printf("locations error: %v", err)
			return fmt.Errorf("Unmarshal error: %v", err)
		}
	} else {
		body, _ := httphandels.HTTPGet(url)
		if err := json.Unmarshal(body, &locations); err != nil {
			fmt.Printf("locations error: %v", err)
			return fmt.Errorf("Unmarshal error: %v", err)
		}
	}

	mapNext = locations.Next
	mapPrevious = locations.Previous
	for i := 0; i < len(locations.Results); i++ {
		fmt.Printf("%s\n", locations.Results[i].Name)
	}
	return nil
}

func commandExplore() error {
	if len(clean) != 1 {
		explore(clean[1])
	} else {
		fmt.Printf("Missing location\n")
	}
	return nil
}

func explore(area string) error {
	url := "https://pokeapi.co/api/v2/location-area/" + area
	var areaPokemon locationArea
	body, ok := cache.Get(url)
	if ok {
		if err := json.Unmarshal(body, &areaPokemon); err != nil {
			fmt.Printf("locations error: %v", err)
			return fmt.Errorf("Unmarshal error: %v", err)
		}
	} else {
		body, _ = httphandels.HTTPGet(url)
		if err := json.Unmarshal(body, &areaPokemon); err != nil {
			fmt.Printf("Incorrect location\n")
			return fmt.Errorf("Unmarshal error: %v", err)
		}
		cache.Add(url, body)
		for _, pokemon := range areaPokemon.PokemonEncounters {
			fmt.Printf("%s\n", pokemon.Pokemon.Name)
		}
	}

	return nil
}

func commandCatch() error {
	if len(clean) != 1 {
		catch(clean[1])
	} else {
		fmt.Printf("Missing Pokemon\n")
	}
	return nil
}

func catch(pokemon string) error {
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)
	var pokemonStats statsPokemon
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemon
	body, ok := cache.Get(url)
	if ok {
		if err := json.Unmarshal(body, &pokemonStats); err != nil {
			fmt.Printf("locations error: %v\n", err)
			return fmt.Errorf("Unmarshal error: %v", err)
		}
	} else {
		body, _ = httphandels.HTTPGet(url)
		if err := json.Unmarshal(body, &pokemonStats); err != nil {
			fmt.Printf("Incorrect Pokemon\n")
			return fmt.Errorf("Unmarshal error: %v", err)
		}
		cache.Add(url, body)
	}
	catchChance := rand.Intn(pokemonStats.BaseExperience)
	if 100-catchChance > 0 {
		fmt.Printf("%s was caught!\n", pokemon)
		pokedex[pokemon] = pokemonStats
	} else {
		fmt.Printf("%s escaped!\n", pokemon)
	}
	return nil
}

func commandInspect() error {
	if len(clean) != 1 {
		inspect(clean[1])
	} else {
		fmt.Printf("Missing Pokemon\n")
	}
	return nil
}

func inspect(pokemon string) error {
	pokeFound := false
	for name, stats := range pokedex {
		if name == pokemon {
			pokeFound = true
			fmt.Printf("Name: %s\n", stats.Name)
			fmt.Printf("Height: %v\n", stats.Height)
			fmt.Printf("Weight: %v\n", stats.Weight)
			fmt.Print("Stats:\n")
			for _, stat := range stats.Stats {
				fmt.Printf("-%s:%v\n", stat.Stat.Name, stat.BaseStat)
			}
			fmt.Print("Types:\n")
			for _, typePoke := range stats.Types {
				fmt.Printf("-%s\n", typePoke.Type.Name)
			}
		}
	}
	if pokeFound == false {
		fmt.Printf("%s not caught yet\n", pokemon)
	}
	return nil
}
