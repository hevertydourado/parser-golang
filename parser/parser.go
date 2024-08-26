package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ParseLogFile parses the Quake 3 Arena log file at the specified filePath.
// It returns a map of matches where each key is a game identifier and the value is a Match struct containing the match details.
func ParseLogFile(filePath string) (map[string]Match, error) {
	file, err := os.Open(filePath) // Open the log file
	if err != nil {
		return nil, err // Return error if the file cannot be opened
	}
	defer file.Close() // Ensure the file is closed when the function exits

	matches := make(map[string]Match) // Map to store all parsed matches
	var currentMatch Match             // Variable to hold the current match being processed
	var matchCount int                 // Counter to track the number of matches

	scanner := bufio.NewScanner(file)  // Scanner to read the file line by line
	killRegex := regexp.MustCompile(`^\s*(\d+:\d+)\sKill:\s(\d+)\s(\d+)\s(\d+):\s(.+)\skilled\s(.+)\sby\s(.+)$`)
	// Regular expression to match kill events in the log file

	for scanner.Scan() { // Iterate over each line in the log file
		line := scanner.Text()

		if strings.Contains(line, "InitGame:") { // Detect the start of a new game
			if matchCount > 0 {
				// If this is not the first match, save the previous match data
				matches[fmt.Sprintf("game_%d", matchCount)] = currentMatch
			}
			matchCount++ // Increment match counter for a new game
			currentMatch = Match{
				Players:      make(map[string]bool),   // Initialize the Players map
				Kills:        make(map[string]int),    // Initialize the Kills map
				KillsByMeans: make(map[string]int),    // Initialize the KillsByMeans map
			}
		} else if strings.Contains(line, "ShutdownGame:") { // Detect the end of a game
			matches[fmt.Sprintf("game_%d", matchCount)] = currentMatch // Save the current match data
		} else if killRegex.MatchString(line) { // Detect a kill event line
			currentMatch = parseKill(line, currentMatch) // Parse the kill event and update the current match data
		}
	}

	// Ensure the last match is stored
	if matchCount > 0 {
		matches[fmt.Sprintf("game_%d", matchCount)] = currentMatch
	}

	if err := scanner.Err(); err != nil { // Check for any errors encountered while reading the file
		return nil, err
	}

	return matches, nil // Return the map of parsed matches
}

// parseKill parses a kill event line from the log and updates the current Match struct accordingly.
func parseKill(line string, match Match) Match {
	killRegex := regexp.MustCompile(`^\s*(\d+:\d+)\sKill:\s(\d+)\s(\d+)\s(\d+):\s(.+)\skilled\s(.+)\sby\s(.+)$`)
	matches := killRegex.FindStringSubmatch(line) // Extract information using the regular expression

	if len(matches) < 8 { // Ensure the line matches the expected format
		fmt.Println("Failed to parse line:", line) // Log an error if the line doesn't match
		return match
	}

	killer := matches[5] // Extract the killer's name
	victim := matches[6] // Extract the victim's name
	means := matches[7]  // Extract the means of killing (e.g., weapon used)

	match.TotalKills++ // Increment the total kill count for the match

	if killer == "<world>" { // Special case where the environment kills the player
		match.Kills[victim]-- // Decrement the victim's kill count
	} else {
		match.Players[killer] = true // Add the killer to the list of players
		match.Players[victim] = true // Add the victim to the list of players
		match.Kills[killer]++        // Increment the killer's kill count
	}

	match.KillsByMeans[means]++ // Increment the count for the specific means of killing

	return match // Return the updated match data
}
