package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func ParseLogFile(filePath string) (map[string]Match, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	matches := make(map[string]Match)
	var currentMatch Match
	var matchCount int

	scanner := bufio.NewScanner(file)
	killRegex := regexp.MustCompile(`^\s*(\d+:\d+)\sKill:\s(\d+)\s(\d+)\s(\d+):\s(.+)\skilled\s(.+)\sby\s(.+)$`)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "InitGame:") {
			if matchCount > 0 {
				matches[fmt.Sprintf("game_%d", matchCount)] = currentMatch
			}
			matchCount++
			currentMatch = Match{
				Players:      make(map[string]bool),
				Kills:        make(map[string]int),
				KillsByMeans: make(map[string]int),
			}
		} else if strings.Contains(line, "ShutdownGame:") {
			matches[fmt.Sprintf("game_%d", matchCount)] = currentMatch
		} else if killRegex.MatchString(line) {
			currentMatch = parseKill(line, currentMatch)
		}
	}

	// Ensure the last match is stored
	if matchCount > 0 {
		matches[fmt.Sprintf("game_%d", matchCount)] = currentMatch
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

func parseKill(line string, match Match) Match {
	killRegex := regexp.MustCompile(`^\s*(\d+:\d+)\sKill:\s(\d+)\s(\d+)\s(\d+):\s(.+)\skilled\s(.+)\sby\s(.+)$`)
	matches := killRegex.FindStringSubmatch(line)

	if len(matches) < 8 {
		fmt.Println("Failed to parse line:", line)
		return match
	}

	killer := matches[5]
	victim := matches[6]
	means := matches[7]

	match.TotalKills++

	if killer == "<world>" {
		match.Kills[victim]--
	} else {
		match.Players[killer] = true
		match.Players[victim] = true
		match.Kills[killer]++
	}

	match.KillsByMeans[means]++

	return match
}
