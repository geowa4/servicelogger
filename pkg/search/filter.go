package search

import (
	"github.com/agnivade/levenshtein"
	"github.com/charmbracelet/bubbles/list"
	"math"
	"slices"
	"strings"
)

type Rank struct {
	contains    bool
	levenshtein int
	rank        list.Rank
}

func FilterFunc(search string, templatesAsStrings []string) []list.Rank {
	ranks := make([]Rank, 0)

	for i, templateAsString := range templatesAsStrings {
		r := Rank{
			contains:    false,
			levenshtein: math.MaxInt,
			rank: list.Rank{
				Index:          i,
				MatchedIndexes: make([]int, 0), // used for underlining matched chars; meh
			},
		}
		if strings.Contains(templateAsString, search) {
			r.contains = true
			r.levenshtein = 0
		} else {
			r.levenshtein = levenshtein.ComputeDistance(search, templateAsString)
		}
		if float64(len(search))/float64(len(templateAsString)) >= 0.1 &&
			!r.contains &&
			float64(r.levenshtein)/float64(len(templateAsString)) >= 0.75 {
			// if the search text is sufficiently large and
			// the search isn't present and
			// the levenshtein distance is relatively large,
			// it's not a match
			continue
		}
		ranks = append(ranks, r)
	}

	// sort by whether the search text is contained by the template,
	// followed by levenshtein distance
	slices.SortFunc(ranks, func(a, b Rank) int {
		if a.contains {
			return -1
		} else if b.contains {
			return 1
		}
		return a.levenshtein - b.levenshtein
	})

	listRanks := make([]list.Rank, 0, len(ranks))
	for _, r := range ranks {
		listRanks = append(listRanks, r.rank)
	}
	return listRanks
}
