package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/lucashthiele/pokedexcli/config"
	"github.com/lucashthiele/pokedexcli/internal/api"
	"github.com/lucashthiele/pokedexcli/internal/pokecache"
	"github.com/lucashthiele/pokedexcli/model"
)

const baseUrl = config.BaseUrl
const cacheDuration = time.Minute * time.Duration(10)
const minCatchChance float64 = 0.05
const maxCatchChance float64 = 0.95
const maxBaseExperience int = 635

type pokedex struct {
	caughtPokemons map[string]model.Pokemon
}

type state struct {
	Previous string
	Next     string
	Cache    *pokecache.Cache
	Pokedex  *pokedex
}

type cliCommand struct {
	name            string
	description     string
	callback        func(*state, []string) error
	validateCommand func(arguments []string) error
}

func printLocations(locations []model.Location) {
	for _, location := range locations {
		fmt.Printf(" - %s\n", location.Name)
	}
}

func printPokemons(pokemons []model.PokemonEncounters) {
	if len(pokemons) > 0 {
		fmt.Printf("Found Pokemon:\n")
	}
	for _, pokemon := range pokemons {
		fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
	}
}

func updateState(state *state, response model.LocationResponse) {
	state.Next = response.Next
	state.Previous = response.Previous
}

func getCachedValue(state *state, url string) []byte {
	if data, found := state.Cache.Get(url); found {
		return data
	}

	return nil
}

func cleanInput(text string) []string {
	slicedStrings := strings.Split(text, " ")

	for i, str := range slicedStrings {
		slicedStrings[i] = strings.ToLower(str)
	}

	return slicedStrings
}

func callbackMap(state *state, params []string) error {
	requestUrl := state.Next
	cachedData := getCachedValue(state, requestUrl)
	if cachedData != nil {
		location, err := api.UnmarshalLocationResponse(cachedData)

		if err != nil {
			return err
		}

		updateState(state, location)
		printLocations(location.Locations)
		return nil
	}

	response, err := api.GetLocations(requestUrl)
	if err != nil {
		return fmt.Errorf("error fetching data from api: %v", err)
	}

	updateState(state, response)
	printLocations(response.Locations)

	data, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("error marshling data from api: %v", err)
	}
	state.Cache.Add(requestUrl, data)

	return nil
}

func callbackMapb(state *state, params []string) error {
	requestUrl := state.Previous
	if requestUrl == "" {
		fmt.Printf("you're on the first page\n")
		return nil
	}

	cachedData := getCachedValue(state, requestUrl)
	if cachedData != nil {
		location, err := api.UnmarshalLocationResponse(cachedData)

		if err != nil {
			return err
		}

		updateState(state, location)
		printLocations(location.Locations)
		return nil
	}

	response, err := api.GetLocations(requestUrl)
	if err != nil {
		return fmt.Errorf("error fetching data from api: %v", err)
	}

	updateState(state, response)
	printLocations(response.Locations)

	data, err := api.MarshalLocation(response)
	if err != nil {
		return err
	}
	state.Cache.Add(requestUrl, data)

	return nil
}

func validateExploreCommand(arguments []string) error {
	if len(arguments) != 2 {
		return fmt.Errorf("expected 1 argument, found %d", len(arguments)-1)
	}
	return nil
}

func callbackExplore(state *state, params []string) error {
	location := params[0]
	if location == "" {
		return fmt.Errorf("you need to provide a location are for this command\n")
	}

	fmt.Printf("exploring %s...\n", location)

	requestUrl, err := url.JoinPath(config.BaseUrl, "location-area", location)
	if err != nil {
		return err
	}

	cachedData := getCachedValue(state, requestUrl)
	if cachedData != nil {
		pokemonResp, err := api.UnmarshalPokemonResponse(cachedData)
		if err != nil {
			return err
		}

		printPokemons(pokemonResp.PokemonEncounters)
		return nil
	}

	response, err := api.GetPokemons(requestUrl)
	if err != nil {
		return fmt.Errorf("error fetching data from api: %v", err)
	}

	printPokemons(response.PokemonEncounters)

	data, err := api.MarshalPokemon(response)
	if err != nil {
		return err
	}
	state.Cache.Add(requestUrl, data)

	return nil
}

func calculateCatchChance(pokemon model.Pokemon) int {
	var difficulty float64 = float64(pokemon.BaseExperience) / float64(maxBaseExperience)

	catchChance := maxCatchChance - (maxCatchChance-minCatchChance)*difficulty

	return int(math.Max(minCatchChance, math.Min(catchChance, maxCatchChance)) * 100)
}

func validateCatchCommand(arguments []string) error {
	if len(arguments) != 2 {
		return fmt.Errorf("expected 1 argument, found %d", len(arguments)-1)
	}
	return nil
}

func callbackCatch(state *state, params []string) error {
	pokemonName := params[0]
	if pokemonName == "" {
		return fmt.Errorf("you need to provide a pokemon for this command\n")
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	requestUrl, err := url.JoinPath(config.BaseUrl, "pokemon", pokemonName)
	if err != nil {
		return err
	}

	pokemon, err := api.GetPokemon(requestUrl)
	if err != nil {
		return fmt.Errorf("error fetching data from api: %v", err)
	}

	catchChance := calculateCatchChance(pokemon) // 1 - 100
	roll := rand.Intn(100)

	if roll < catchChance {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		state.Pokedex.caughtPokemons[pokemon.Name] = pokemon
		return nil
	}

	fmt.Printf("%s escaped!\n", pokemon.Name)
	return nil
}

func callbackHelp(state *state, params []string) error {
	commands := getCommandMap()
	fmt.Printf("\nWelcome to the Pokedex!\nUsage:\n\n")

	for _, value := range commands {
		fmt.Printf("%s: %s\n", value.name, value.description)
	}
	fmt.Print("\n")

	return nil
}

func callbackExit(state *state, params []string) error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func getCommandMap() map[string]cliCommand {
	return map[string]cliCommand{
		"map": {
			name:            "map",
			description:     "Displays the name of 20 location areas in Pokemon world. Any subsequent call to map, will return the next 20 locations",
			callback:        callbackMap,
			validateCommand: func(arguments []string) error { return nil },
		},
		"mapb": {
			name:            "mapb",
			description:     "Displays the 20 previous locations areas",
			callback:        callbackMapb,
			validateCommand: func(arguments []string) error { return nil },
		},
		"explore": {
			name:            "explore <location-name>",
			description:     "Display all pokemons within that location area",
			callback:        callbackExplore,
			validateCommand: validateExploreCommand,
		},
		"catch": {
			name:            "catch <pokemon-name>",
			description:     "Try to catch the given Pokémon. The chance of catching it is calculated based on its base experience",
			callback:        callbackCatch,
			validateCommand: validateCatchCommand,
		},
		"help": {
			name:            "help",
			description:     "Displays a help message",
			callback:        callbackHelp,
			validateCommand: func(arguments []string) error { return nil },
		},
		"exit": {
			name:            "exit",
			description:     "Exits the program",
			callback:        callbackExit,
			validateCommand: func(arguments []string) error { return nil },
		},
	}
}

func main() {
	commandMap := getCommandMap()

	cache := pokecache.NewCache(cacheDuration)

	state := &state{
		Previous: "",
		Next:     baseUrl + "/location-area",
		Cache:    &cache,
		Pokedex: &pokedex{
			caughtPokemons: make(map[string]model.Pokemon),
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("Pokedex > ")
		var option string
		for scanner.Scan() {
			option = scanner.Text()
			break
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}

		arguments := cleanInput(option)
		command, found := commandMap[arguments[0]]

		if !found {
			fmt.Fprintln(os.Stderr, "Unknow Command")
			callbackHelp(nil, nil)
			continue
		}

		err := command.validateCommand(arguments)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			continue
		}

		err = command.callback(state, arguments[1:])
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}
	}
}
