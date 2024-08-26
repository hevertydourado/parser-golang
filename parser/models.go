package parser

// Match represents the data structure for a single game match.
type Match struct {
	TotalKills   int               // Total number of kills that occurred in the match
	Players      map[string]bool   // Set of players who participated in the match
	Kills        map[string]int    // Map of player names to the number of kills they achieved
	KillsByMeans map[string]int    // Map of kill means (e.g., weapons used) to the number of kills associated with each
}
