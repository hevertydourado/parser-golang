package main

import (
        "encoding/json"
        "fmt"
        "os"
        "quake-log-parser/parser"
        "sort"
        "strconv"
)

type Game struct {
        TotalKills   int               // Total number of kills in the game
        Players      map[string]bool   // Set of unique players in the game
        Kills        map[string]int    // Map of player names to their respective kill counts
        KillsByMeans map[string]int    // Map of kill means (e.g., weapon type) to their respective counts
}

type ReportGenerator struct {
        Games map[int]Game             // Map of game numbers to their corresponding Game data
}

// generateMatchReports generates a report for each match, including total kills, players, and individual player kills.
func (rg *ReportGenerator) generateMatchReports() map[string]interface{} {
        matchReports := make(map[string]interface{})
        for gameNumber, game := range rg.Games {
                gameKey := fmt.Sprintf("game_%02d", gameNumber) // Format game number as "game_XX"
                matchReports[gameKey] = map[string]interface{}{
                        "total_kills": game.TotalKills,       // Total kills in the match
                        "players":     getPlayersList(game.Players), // List of players in the match
                        "kills":       game.Kills,            // Kill counts per player
                }
        }
        return matchReports
}

// generateKillByMeansReport generates a report for each match showing the number of kills by different means (e.g., weapons).
func (rg *ReportGenerator) generateKillByMeansReport() map[string]interface{} {
        meansReports := make(map[string]interface{})
        for gameNumber, game := range rg.Games {
                gameKey := fmt.Sprintf("game_%02d", gameNumber) // Format game number as "game_XX"
                meansReports[gameKey] = map[string]interface{}{
                        "kills_by_means": game.KillsByMeans, // Kills categorized by means
                }
        }
        return meansReports
}

// generatePlayerRanking generates a ranking of players based on their total kills across all games.
func (rg *ReportGenerator) generatePlayerRanking() []map[string]interface{} {
        playerScores := make(map[string]int)
        for _, game := range rg.Games {
                for player, kills := range game.Kills {
                        playerScores[player] += kills // Accumulate kills for each player
                }
        }

        // Sort players by total kills in descending order
        sortedPlayers := make([]map[string]interface{}, 0, len(playerScores))
        for player, kills := range playerScores {
                sortedPlayers = append(sortedPlayers, map[string]interface{}{
                        "player": player,
                        "kills":  kills,
                })
        }
        sort.Slice(sortedPlayers, func(i, j int) bool {
                return sortedPlayers[i]["kills"].(int) > sortedPlayers[j]["kills"].(int) // Sort by kills descending
        })

        return sortedPlayers
}

// saveReportsAsJSON saves the generated reports as a JSON file with proper formatting.
func (rg *ReportGenerator) saveReportsAsJSON(matchReports, meansReports map[string]interface{}, playerRanking []map[string]interface{}, filename string) error {
        outputData := map[string]interface{}{
                "match_reports":        matchReports,        // Reports for each match
                "kill_by_means_reports": meansReports,        // Reports for kills by means
                "player_ranking":       playerRanking,        // Player ranking report
        }

        file, err := os.Create(filename) // Create a new file to save the JSON data
        if err != nil {
                return err
        }
        defer file.Close() // Ensure the file is closed after writing

        encoder := json.NewEncoder(file)
        encoder.SetIndent("", "  ") // Format JSON with indentation for readability
        return encoder.Encode(outputData) // Write the JSON data to the file
}

// getPlayersList converts the map of players to a sorted list of player names.
func getPlayersList(players map[string]bool) []string {
        playerList := make([]string, 0, len(players))
        for player := range players {
                playerList = append(playerList, player)
        }
        sort.Strings(playerList) // Sort the players alphabetically for consistency
        return playerList
}

func main() {
        // Define the path to the log file to be parsed
        logFilePath := "logs/quake_game.log"
        
        // Parse the log file to extract match data
        matches, err := parser.ParseLogFile(logFilePath)
        if err != nil {
                fmt.Println("Error parsing log file:", err)
                return
        }

        // Convert parsed matches into Game structures
        games := make(map[int]Game)
        for matchNumber, matchData := range matches {
                gameNumber, err := strconv.Atoi(matchNumber[5:]) // Extract the game number from the match identifier
                if err != nil {
                        fmt.Println("Error parsing game number:", err)
                        continue
                }
                game := Game{
                        TotalKills:   matchData.TotalKills,   // Total kills recorded in the match
                        Players:      matchData.Players,      // Players who participated in the match
                        Kills:        matchData.Kills,        // Kill counts per player
                        KillsByMeans: matchData.KillsByMeans, // Kill counts categorized by means (e.g., weapons)
                }
                games[gameNumber] = game
        }

        // Instantiate a ReportGenerator with the parsed game data
        reportGen := ReportGenerator{Games: games}

        // Generate the match reports, kills by means reports, and player ranking
        matchReports := reportGen.generateMatchReports()
        meansReports := reportGen.generateKillByMeansReport()
        playerRanking := reportGen.generatePlayerRanking()

        // Save the generated reports to a JSON file
        err = reportGen.saveReportsAsJSON(matchReports, meansReports, playerRanking, "quake_report.json")
        if err != nil {
                fmt.Println("Error saving JSON report:", err)
                return
        }

        fmt.Println("Report generated successfully in quake_report.json") // Confirm successful report generation
}
