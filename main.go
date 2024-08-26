package main

import (
	"encoding/json"
	"fmt"
	"os"
	"quake-log-parser/parser" // Aqui você deve importar o seu módulo de parser
	"sort"
	"strconv"
)

type Game struct {
	TotalKills   int
	Players      map[string]bool
	Kills        map[string]int
	KillsByMeans map[string]int
}

type ReportGenerator struct {
	Games map[int]Game
}

func (rg *ReportGenerator) generateMatchReports() map[string]interface{} {
	matchReports := make(map[string]interface{})
	for gameNumber, game := range rg.Games {
		gameKey := fmt.Sprintf("game_%02d", gameNumber)
		matchReports[gameKey] = map[string]interface{}{
			"total_kills": game.TotalKills,
			"players":     getPlayersList(game.Players),
			"kills":       game.Kills,
		}
	}
	return matchReports
}

func (rg *ReportGenerator) generateKillByMeansReport() map[string]interface{} {
	meansReports := make(map[string]interface{})
	for gameNumber, game := range rg.Games {
		gameKey := fmt.Sprintf("game_%02d", gameNumber)
		meansReports[gameKey] = map[string]interface{}{
			"kills_by_means": game.KillsByMeans,
		}
	}
	return meansReports
}

func (rg *ReportGenerator) generatePlayerRanking() []map[string]interface{} {
	playerScores := make(map[string]int)
	for _, game := range rg.Games {
		for player, kills := range game.Kills {
			playerScores[player] += kills
		}
	}

	// Ordena os jogadores pelo total de kills em ordem decrescente
	sortedPlayers := make([]map[string]interface{}, 0, len(playerScores))
	for player, kills := range playerScores {
		sortedPlayers = append(sortedPlayers, map[string]interface{}{
			"player": player,
			"kills":  kills,
		})
	}
	sort.Slice(sortedPlayers, func(i, j int) bool {
		return sortedPlayers[i]["kills"].(int) > sortedPlayers[j]["kills"].(int)
	})

	return sortedPlayers
}

func (rg *ReportGenerator) saveReportsAsJSON(matchReports, meansReports map[string]interface{}, playerRanking []map[string]interface{}, filename string) error {
	outputData := map[string]interface{}{
		"match_reports":        matchReports,
		"kill_by_means_reports": meansReports,
		"player_ranking":       playerRanking,
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(outputData)
}

func getPlayersList(players map[string]bool) []string {
	playerList := make([]string, 0, len(players))
	for player := range players {
		playerList = append(playerList, player)
	}
	sort.Strings(playerList) // Ordena os jogadores alfabeticamente para consistência
	return playerList
}

func main() {
	// Parse the log file to get the games data
	logFilePath := "logs/quake_game.log" // Atualize o caminho do arquivo de log conforme necessário
	matches, err := parser.ParseLogFile(logFilePath)
	if err != nil {
		fmt.Println("Error parsing log file:", err)
		return
	}

	// Convert parsed matches to the Game structure
	games := make(map[int]Game)
	for matchNumber, matchData := range matches {
		gameNumber, err := strconv.Atoi(matchNumber[5:]) // Extrair o número do jogo do nome, assumindo que o nome está no formato "game_X"
		if err != nil {
			fmt.Println("Error parsing game number:", err)
			continue
		}
		game := Game{
			TotalKills:   matchData.TotalKills,
			Players:      matchData.Players,
			Kills:        matchData.Kills,
			KillsByMeans: matchData.KillsByMeans,
		}
		games[gameNumber] = game
	}

	reportGen := ReportGenerator{Games: games}

	matchReports := reportGen.generateMatchReports()
	meansReports := reportGen.generateKillByMeansReport()
	playerRanking := reportGen.generatePlayerRanking()

	err = reportGen.saveReportsAsJSON(matchReports, meansReports, playerRanking, "quake_report.json")
	if err != nil {
		fmt.Println("Error saving JSON report:", err)
		return
	}

	fmt.Println("Report generated successfully in quake_report.json")
}
