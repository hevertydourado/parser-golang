package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

type MatchReport struct {
	TotalKills   int            `json:"total_kills"`
	Players      []string       `json:"players"`
	Kills        map[string]int `json:"kills"`
}

type Output struct {
	Reports       map[string]MatchReport `json:"reports"`
	PlayerRanking []map[string]interface{} `json:"player_ranking"`
}

func main() {
	// Read the JSON file
	data, err := ioutil.ReadFile("quake_report.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	// Parse the JSON into our struct
	var output Output
	if err := json.Unmarshal(data, &output); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Debug: Check if we have match reports
	if len(output.Reports) == 0 {
		fmt.Println("No match reports found in JSON.")
		return
	}

	// Ask user how many games to display
	fmt.Print("Enter the number of games to display (default is 10): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	numGames := 10
	if input != "" {
		if n, err := strconv.Atoi(input); err == nil && n > 0 {
			numGames = n
		} else {
			fmt.Println("Invalid input. Using default value of 10.")
		}
	}

	// Extract and sort the game keys
	gameKeys := make([]string, 0, len(output.Reports))
	for game := range output.Reports {
		gameKeys = append(gameKeys, game)
	}
	sort.Strings(gameKeys)

	// Display information for the first `numGames` matches
	for i, game := range gameKeys {
		if i >= numGames {
			break
		}
		report := output.Reports[game]
		fmt.Printf("Game: %s\n", game)
		fmt.Printf("  Total Kills: %d\n", report.TotalKills)
		fmt.Printf("  Players: %v\n", report.Players)
		fmt.Println("  Kills:")
		for player, kills := range report.Kills {
			fmt.Printf("    %s: %d\n", player, kills)
		}
		fmt.Println()
	}
}
