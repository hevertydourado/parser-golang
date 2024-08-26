package main

import (
    "encoding/json"
    "fmt"
    "os"
    "quake-log-parser/parser"
    "sort"
    "strconv"
)

// Game represents the structure holding the game data.
type Game struct {
    TotalKills   int            // TotalKills is the total number of kills in a game.
    Players      map[string]bool // Players is a set of unique player names in the game.
    Kills        map[string]int // Kills maps each player to the number of kills they have.
    KillsByMeans map[string]int // KillsByMeans maps the kill methods to their corresponding counts.
}

// ReportGenerator is responsible for generating various reports from the game data.
type ReportGenerator struct {
    Games map[int]Game // Games is a map of all games indexed by their game number.
}

// generateMatchReports generates a detailed report of each match.
// It returns a map where the key is the game identifier and the value is the match details.
func (rg *ReportGenerator) generateMatchReports() map[string]interface{} {
    matchReports := make(map[string]interface{})
    for gameNumber, game := range rg.Games {
        gameKey := fmt.Sprintf("game_%02d", gameNumber) // Formats the game number as "game_XX".
        matchReports[gameKey] = map[string]interface{}{
            "total_kills": game.TotalKills,     // Total number of kills in the game.
            "players":     getPlayersList(game.Players), // List of players in the game.
            "kills":       game.Kills,          // Mapping of players to their kill counts.
        }
    }
    return matchReports
}

// generateKillByMeansReport generates a report of kills by the means used (e.g., weapons).
// It returns a map where the key is the game identifier and the value is the kill methods breakdown.
func (rg *ReportGenerator) generateKillByMeansReport() map[string]interface{} {
    meansReports := make(map[string]interface{})
    for gameNumber, game := range rg.Games {
        gameKey := fmt.Sprintf("game_%02d", gameNumber) // Formats the game number as "game_XX".
        meansReports[gameKey] = map[string]interface{}{
            "kills_by_means": game.KillsByMeans, // Mapping of kill methods to their counts.
        }
    }
    return meansReports
}

// getPlayersList converts the set of players (a map) to a sorted slice of player names.
func getPlayersList(players map[string]bool) []string {
    playerList := make([]string, 0, len(players))
    for player := range players {
        playerList = append(playerList, player) // Append each player name to the list.
    }
    sort.Strings(playerList) // Sort the player names alphabetically.
    return playerList
}

// parseGameLog processes the game log file, extracting game data and populating the report generator.
func parseGameLog(filepath string) (*ReportGenerator, error) {
    logParser, err := parser.NewParser(filepath)
    if err != nil {
        return nil, fmt.Errorf("failed to create log parser: %w", err) // Error handling if parser creation fails.
    }

    rg := &ReportGenerator{
        Games: make(map[int]Game),
    }

    gameNumber := 0
    for logParser.Scan() {
        event := logParser.Event()

        // If a new game starts, increment the game number and initialize the game data structure.
        if event.Type == parser.EventTypeInitGame {
            gameNumber++
            rg.Games[gameNumber] = Game{
                Players:      make(map[string]bool),
                Kills:        make(map[string]int),
                KillsByMeans: make(map[string]int),
            }
        }

        // Handle different event types related to kills and player connections.
        game := rg.Games[gameNumber]
        switch event.Type {
        case parser.EventTypeKill:
            game.TotalKills++
            game.Kills[event.Killer]++
            game.KillsByMeans[event.MeanOfDeath]++
        case parser.EventTypeClientUserinfoChanged:
            game.Players[event.PlayerName] = true
        }

        rg.Games[gameNumber] = game
    }

    if err := logParser.Err(); err != nil {
        return nil, fmt.Errorf("error occurred during log parsing: %w", err) // Error handling if log parsing fails.
    }

    return rg, nil
}

// main function is the entry point of the application.
// It reads the log file path from command line arguments, parses the log, generates reports, and outputs them as JSON.
func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: quake-log-parser <logfile>")
        os.Exit(1)
    }

    logfile := os.Args[1]
    rg, err := parseGameLog(logfile)
    if err != nil {
        fmt.Printf("Error parsing game log: %v\n", err)
        os.Exit(1)
    }

    matchReports := rg.generateMatchReports()
    meansReports := rg.generateKillByMeansReport()

    combinedReports := map[string]interface{}{
        "match_reports":      matchReports,      // Detailed match reports.
        "kill_by_means_report": meansReports,    // Kills by means reports.
    }

    jsonData, err := json.MarshalIndent(combinedReports, "", "  ")
    if err != nil {
        fmt.Printf("Error generating JSON: %v\n", err)
        os.Exit(1)
    }

    fmt.Println(string(jsonData)) // Output the generated JSON to stdout.
}
