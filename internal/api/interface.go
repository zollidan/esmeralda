package api

type MatchFetcher interface {
	GetMatches(filter MatchesFilter) ([]Match, int, error)
}