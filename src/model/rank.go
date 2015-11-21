package model

import (
	"fmt"
)

func (st *StackedTranslation) GetVotes() []*Vote {
	// entry := st.Entry.Entries[0]
	results := query("select "+voteFields+" from Votes where TranslationID = ?", st.Parts[0].ID()).rows(parseVote)

	votes := make([]*Vote, len(results))
	for i, result := range results {
		if vote, ok := result.(Vote); ok {
			votes[i] = &vote
		}
	}
	return votes
}

func GetPreferredTranslations(language string) []*StackedTranslation {
	entries := stackEntries(GetEntries())
	pref := make([]*StackedTranslation, 0, len(entries))
	for _, entry := range entries {
		translations := entry.GetTranslations(language)
		selected := PickPreferredTranslation(entry.RankTranslations(translations, false))
		if selected != nil {
			pref = append(pref, selected)
		}
	}

	return pref
}

type RankTranslation struct {
	Translation *StackedTranslation
	Rank        int
}

func PickPreferredTranslation(translations []RankTranslation) *StackedTranslation {
	if len(translations) == 0 {
		return nil
	}

	for _, tr := range translations {
		if tr.Translation.IsPreferred {
			return tr.Translation
		}
	}
	return translations[0].Translation
}

func (entry *StackedEntry) RankTranslations(translations []*StackedTranslation, save bool) []RankTranslation {
	if len(translations) == 0 {
		return nil
	}

	language := translations[0].Language
	lead := GetLanguageLead(language)

	ln := len(translations)

	// count votes
	scores := make(map[string]int, ln)
	upvoters := make(map[string][]string, ln)
	downvoters := make(map[string][]string, ln)

	for _, translation := range translations {
		scores[translation.FullText] = 0
		upvoters[translation.FullText] = make([]string, 0, ln)
		downvoters[translation.FullText] = make([]string, 0, ln)
	}

	for _, translation := range translations {
		upvoters[translation.FullText] = append(upvoters[translation.FullText], translation.Translator)

		for _, vote := range translation.GetVotes() {
			if vote.Vote {
				upvoters[translation.FullText] = append(upvoters[translation.FullText], vote.Voter.Email)
			} else {
				downvoters[translation.FullText] = append(downvoters[translation.FullText], vote.Voter.Email)
			}
		}
	}

	for text, ups := range upvoters {
		for _, voter := range ups {
			voteWeight := 2
			if lead != nil && voter == lead.Email {
				voteWeight++;
			}

			scores[text] += voteWeight
		}

		for _, voter := range downvoters[text] {
			voteWeight := 2
			if lead != nil && voter == lead.Email {
				voteWeight++;
			}

			scores[text] -= voteWeight
		}
	}

	if Debug >= 2 {
		fmt.Println(" - Voting scores:", scores)
	}

	// get translations from people who haven't upvoted
	// for _, translation := range translations {
	// 	if !hasUpvoted[translation.Translator] {
	// 		voteWeight := 2
	// 		if lead != nil && translation.Translator == lead.Email {
	// 			voteWeight++;
	// 		}
	// 		scores[translation.FullText] += voteWeight
	// 	}
	// }

	// find the highest rank
	topScore := 0
	topScoringText := ""
	for text, score := range scores {
		if score > topScore {
			topScore = score
			topScoringText = text
		}
	}

	// check if more than one translation has near that score
	threshold := topScore - 1
	numNearTopRank := 0
	for _, score := range scores {
		if score >= threshold {
			numNearTopRank++
		}
	}
	isConflicted := numNearTopRank > 1
	if isConflicted && Debug >= 1 {
		fmt.Println("Conflict!", numNearTopRank, "translations for:", entry.FullText)
	}

	// update their flags
	for _, translation := range translations {
		translation.IsConflicted = isConflicted && scores[translation.FullText] >= threshold
		translation.IsPreferred = !isConflicted && translation.FullText == topScoringText

		for _, part := range translation.Parts {
			part.IsConflicted = translation.IsConflicted
			part.IsPreferred = translation.IsPreferred
			if save {
				if Debug >= 2 {
					fmt.Println(" - Saving translation part:", part)
				}
				part.Save(false)
			}
		}
	}

	if Debug >= 2 {
		fmt.Println(" - Ranking for return")
	}

	// return the marked and ranked translations
	ranked := make([]RankTranslation, len(translations))
	for i, translation := range translations {
		score := scores[translation.FullText]
		ranked[i] = RankTranslation{translation, score}
	}
	return ranked
}


func (se *StackedEntry) MarkConflicts(language string) {
	if Debug >= 1 {
		fmt.Println("Marking conflicts in '"+se.FullText+"' ("+language+")")
	}
	translations := se.GetTranslations(language)
	if Debug >= 1 {
		fmt.Println("Ranking", len(translations), "translations")
	}
	se.RankTranslations(translations, true)
}


func MarkAllConflicts() {
	stackedEntries := GetStackedEntries("", "", "", "", "", "", nil)
	fmt.Println("Loaded", len(stackedEntries), "stacked entries")

	for _, se := range stackedEntries {
		for _, lang := range Languages {
			se.MarkConflicts(lang)
		}
	}
}