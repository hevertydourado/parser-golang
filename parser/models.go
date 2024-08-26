package parser

type Match struct {
	TotalKills   int
	Players      map[string]bool
	Kills        map[string]int
	KillsByMeans map[string]int
}
