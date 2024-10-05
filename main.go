package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/lucashthiele/pokedexcli/internal/api"
	"github.com/lucashthiele/pokedexcli/model"
)

const baseUrl = "https://pokeapi.co/api/v2"

type config struct {
	Previous string
	Next     string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func main() {
	commandMap := getCommandMap()

	config := &config{
		Previous: "",
		Next:     baseUrl + "/location-area",
	}

	for {
		fmt.Printf("pokedex > ")
		scanner := bufio.NewScanner(os.Stdin)
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
	fmt.Printf("\nWelcome to the pokedex\nUsage:\n\n")

	for _, value := range commands {
		fmt.Printf("%s: %s\n", value.name, value.description)
	}
	fmt.Print("\n")

	return nil
}

func callbackExit(config *config) error {
	os.Exit(0)
	return nil
}

func callbackMap(config *config) error {
	response, err := api.GetLocations(config.Next)
	if err != nil {
		return fmt.Errorf("error fetching data from api: %v", err)
	}

	updateConfig(config, response)

	printLocations(response.Locations)

	return nil
}

func callbackMapb(config *config) error {
	response, err := api.GetLocations(config.Previous)
	if err != nil {
		return fmt.Errorf("error fetching data from api: %v", err)
	}

	updateConfig(config, response)

	printLocations(response.Locations)

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
