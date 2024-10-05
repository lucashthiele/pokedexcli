package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func main() {
	commandMap := getCommandMap()

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
			callbackHelp()
		} else {
			command.callback()
		}
	}
}

func getCommandMap() map[string]cliCommand {
	return map[string]cliCommand{
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

func callbackHelp() error {
	commands := getCommandMap()
	fmt.Printf("\nWelcome to the pokedex\nUsage:\n\n")

	for _, value := range commands {
		fmt.Printf("%s: %s\n", value.name, value.description)
	}
	fmt.Print("\n")

	return nil
}

func callbackExit() error {
	os.Exit(0)
	return nil
}
