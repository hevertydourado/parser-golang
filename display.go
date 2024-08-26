package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "sort"
    "strings"
)

// Report represents the entire structure of the parsed JSON data.
type Report struct {
    KillByMeansReports map[string]KillByMeansReport `json:"kill_by_means_reports"`
    MatchReports       map[string]MatchReport       `json:"match_reports"`
    PlayerRanking      []PlayerRanking              `json:"player_ranking"`
}

// KillByMeansReport represents the kills by means for a specific game.
type KillByMeansReport struct {
    KillsByMeans map[string]int `json:"kills_by_means"`
}

// MatchReport represents the report of a specific match.
type MatchReport struct {
    TotalKills int               `json:"total_kills"`
    Players    []string          `json:"players"`
    Kills      map[string]int    `json:"kills"`
}

// PlayerRanking represents the ranking of players based on their kills.
type PlayerRanking struct {
    Player string `json:"player"`
    Kills  int    `json:"kills"`
}

// main is the entry point of the program.
func main() {
    // Load the JSON file containing the report data.
    data, err := ioutil.ReadFile("quake_report.json")
    if err != nil {
        log.Fatalf("Failed to read JSON file: %s", err)
    }

    // Parse the JSON data into the Report structure.
    var report Report
    err = json.Unmarshal(data, &report)
    if err != nil {
        log.Fatalf("Failed to parse JSON data: %s", err)
    }

    // Display the player ranking.
    fmt.Println("===== Player Ranking =====")
    displayPlayerRanking(report.PlayerRanking)

    // Display the match reports.
    fmt.Println("\n===== Match Reports =====")
    for match, matchReport := range report.MatchReports {
        displayMatchReport(match, matchReport)
    }

    // Display the kill by means reports.
    fmt.Println("\n===== Kill by Means Reports =====")
    for game, killByMeansReport := range report.KillByMeansReports {
        displayKillByMeansReport(game, killByMeansReport)
    }
}

// displayPlayerRanking prints the ranking of players in a formatted way.
func displayPlayerRanking(rankings []PlayerRanking) {
    // Sort players by the number of kills in descending order.
    sort.Slice(rankings, func(i, j int) bool {
        return rankings[i].Kills > rankings[j].Kills
    })

    // Print the player ranking.
    for i, player := range rankings {
        fmt.Printf("%d. %s - %d kills\n", i+1, player.Player, player.Kills)
    }
}

// displayMatchReport prints the details of a specific match.
func displayMatchReport(match string, report MatchReport) {
    fmt.Printf("\nMatch: %s\n", match)
    fmt.Printf("Total Kills: %d\n", report.TotalKills)
    fmt.Printf("Players: %s\n", strings.Join(report.Players, ", "))

    fmt.Println("Kills:")
    // Sort players by kills in descending order for better readability.
    sortedPlayers := sortKills(report.Kills)
    for _, player := range sortedPlayers {
        fmt.Printf("  %s: %d kills\n", player.Name, player.Kills)
    }
}

// displayKillByMeansReport prints the kills by means for a specific game.
func displayKillByMeansReport(game string, report KillByMeansReport) {
    fmt.Printf("\nGame: %s\n", game)
    fmt.Println("Kills by Means:")

    // Sort the kill by means data by number of kills in descending order.
    sortedMeans := sortKills(report.KillsByMeans)
    for _, mean := range sortedMeans {
        fmt.Printf("  %s: %d kills\n", mean.Name, mean.Kills)
    }
}

// sortKills sorts a map of kills by the number of kills in descending order.
func sortKills(kills map[string]int) []PlayerKill {
    var sortedKills []PlayerKill
    for name, kills := range kills {
        sortedKills = append(sortedKills, PlayerKill{Name: name, Kills: kills})
    }
    sort.Slice(sortedKills, func(i, j int) bool {
        return sortedKills[i].Kills > sortedKills[j].Kills
    })
    return sortedKills
}

// PlayerKill is a helper struct to hold player names and kill counts together.
type PlayerKill struct {
    Name  string
    Kills int
}
