package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lucashthiele/pokedexcli/internal/api"
	"github.com/lucashthiele/pokedexcli/internal/logs"
	"github.com/lucashthiele/pokedexcli/internal/pokecache"
	"github.com/lucashthiele/pokedexcli/model"
)

const baseUrl = "https://pokeapi.co/api/v2"
const cacheDuration = time.Second * time.Duration(10)

type config struct {
	Previous string
	Next     string
	Cache    *pokecache.Cache
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func getCommandMap() map[string]cliCommand {
	return map[string]cliCommand{
		"map": {
			name:        "map",
			description: "Displays the name of 20 location areas in Pokemon world. Any subsequent call to map, will return the next 20 locations",
			callback:    callbackMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the 20 previous locations areas",
			callback:    callbackMapb,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    callbackHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exits the program",
			callback:    callbackExit,
		},
	}
}

func callbackHelp(config *config) error {
	commands := getCommandMap()
	fmt.Printf("\nWelcome to the Pokedex!\nUsage:\n\n")

	for _, value := range commands {
		fmt.Printf("%s: %s\n", value.name, value.description)
	}
	fmt.Print("\n")

	return nil
}

func callbackExit(config *config) error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func checkCache(config *config, url string) (bool, error) {
	if data, found := config.Cache.Get(url); found {
		logs.Log("found entry in cache")
		location := model.LocationResponse{}

		err := json.Unmarshal(data, &location)
		if err != nil {
			return false, err
		}

		updateConfig(config, location)

		printLocations(location.Locations)

		return true, nil
	}

	logs.Log("entry not found in cache, sending request")

	return false, nil
}

func callbackMap(config *config) error {
	requestUrl := config.Next
	foundCache, err := checkCache(config, requestUrl)
	if foundCache {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error fetching data from cache: %v", err)
	}

	response, err := api.GetLocations(requestUrl)
	if err != nil {
		return fmt.Errorf("error fetching data from api: %v", err)
	}

	updateConfig(config, response)

	printLocations(response.Locations)

	data, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("error marshling data from api: %v", err)
	}
	config.Cache.Add(requestUrl, data)

	return nil
}

func callbackMapb(config *config) error {
	requestUrl := config.Previous
	if requestUrl == "" {
		fmt.Printf("you're on the first page\n")
		return nil
	}

	foundCache, err := checkCache(config, requestUrl)
	if foundCache {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error fetching data from cache: %v", err)
	}

	response, err := api.GetLocations(requestUrl)
	if err != nil {
		return fmt.Errorf("error fetching data from api: %v", err)
	}

	updateConfig(config, response)

	printLocations(response.Locations)

	data, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("error marshling data from api: %v", err)
	}
	config.Cache.Add(requestUrl, data)

	return nil
}

func printLocations(locations []model.Location) {
	for _, location := range locations {
		fmt.Printf("%s\n", location.Name)
	}
}

func updateConfig(config *config, response model.LocationResponse) {
	config.Next = response.Next
	config.Previous = response.Previous
}

func cleanInput(text string) []string {
	slicedStrings := strings.Split(text, " ")

	for i, str := range slicedStrings {
		slicedStrings[i] = strings.ToLower(str)
	}

	return slicedStrings
}

func main() {
	commandMap := getCommandMap()

	cache := pokecache.NewCache(cacheDuration)

	config := &config{
		Previous: "",
		Next:     baseUrl + "/location-area",
		Cache:    &cache,
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

		command, ok := commandMap[option]
		if !ok {
			fmt.Fprintln(os.Stderr, "Unknow Command")
			callbackHelp(nil)
		} else {
			command.callback(config)
		}
	}
}
