# PokedexCLI

PokedexCLI is a command-line interface (CLI) application written in Go (Golang) that simulates a Pokédex. It allows users to search for Pokémon names, types, and stats. The app runs as a REPL (Read-Eval-Print Loop), meaning you can continuously input commands to interact with the Pokédex.

## What is a Pokédex?

A **Pokédex** is an in-universe device from the Pokémon franchise used to catalog and provide information about various Pokémon species. In this CLI version, you can search for Pokémon by name, view their types, and see relevant stats like health, attack, defense, etc.

## Requirements

- **Go 1.22.5**

## Installation

### 1. Clone the repository

To get started, clone the repository using Git:

```bash
git clone https://github.com/lucashthiele/pokedexcli.git
```

### 2. Build the application

Navigate to the cloned repository directory and build the application:

```bash
cd pokedexcli
go build
```

This will compile the Go code and produce an executable file.

### 3. Start the REPL

Once built, you can start the CLI REPL by running the following command:

```bash
./pokedexcli
```

You will enter an interactive mode where you can begin searching for Pokémon by their names or types, and view their stats.

## Commands

The following commands are available in the PokedexCLI:

- **map**:  
  Displays the names of 20 location areas in the Pokémon world. If called multiple times, each call will return the next 20 locations.

- **mapb**:  
  Displays the 20 previous location areas in the Pokémon world, allowing you to navigate back through the list.

- **help**:  
  Displays a help message, providing details on how to use the CLI and its commands.

- **exit**:  
  Exits the application, ending the REPL session.
  
## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.