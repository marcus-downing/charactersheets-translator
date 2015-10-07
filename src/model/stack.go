package model

import (
	// "code.google.com/p/go.crypto/bcrypt"
	"fmt"
	// "math/rand"
	"sort"
	"strconv"
	"strings"
	// "time"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

type StackedEntry struct {
	FullText     string
	Entries      []*Entry
	EntrySources []*EntrySource
	Count        int
	SourceCount  int
}

func (se *StackedEntry) ID() uint64 {
	return se.Entries[0].ID()
}

func stackEntries(entries []*Entry) []*StackedEntry {
	fmt.Println("Stacking", len(entries), "entries")
	stacks := make(map[string][]*Entry, len(entries))
	unstacked := make([]*Entry, 0, len(entries))
	for _, entry := range entries {
		if entry.PartOf != "" {
			if stacks[entry.PartOf] == nil {
				stacks[entry.PartOf] = make([]*Entry, 0, 10)
			}
			stacks[entry.PartOf] = append(stacks[entry.PartOf], entry)
		} else {
			unstacked = append(unstacked, entry)
		}
	}

	// put entries in order
	if Debug >= 1 {
		fmt.Println("Sorting enties")
	}
	values := make([]*StackedEntry, 0, len(stacks)+len(unstacked))
	for _, stack := range stacks {
		sort.Sort(entriesByIndex(stack))
		values = append(values, &StackedEntry{
			FullText: stack[0].PartOf,
			Entries:  stack,
		})
	}
	for _, entry := range unstacked {
		values = append(values, &StackedEntry{
			FullText: entry.Original,
			Entries:  []*Entry{entry},
		})
	}

	// load sources
	if Debug >= 1 {
		fmt.Println("Loading sources")
	}
	sources := make(map[uint64]*Source, 500)
	for _, source := range GetSources() {
		sources[source.ID()] = source
	}

	// calculate totals
	if Debug >= 1 {
		fmt.Println("Calculating totals")
	}
	for _, se := range values {
		entrySources := make(map[uint64]*EntrySource, len(se.Entries)*10)
		for _, entry := range se.Entries {
			for _, placeholder := range GetSourceIDsForEntry(entry) {
				if source, ok := sources[placeholder.SourceID]; ok {
					entrySources[placeholder.SourceID] = &EntrySource{*entry, *source, placeholder.Count}
				}
			}
		}
		count := 0
		esv := make([]*EntrySource, 0, len(entrySources))
		for _, es := range entrySources {
			esv = append(esv, es)
			count += es.Count
		}
		se.EntrySources = esv
		se.Count = count
	}
	return values
}

func sortStacks(values []*StackedEntry, sortBy string) []*StackedEntry {
	if Debug >= 1 {
		fmt.Println("Sorting stacks by:", sortBy)
	}
	switch sortBy {
	case "", "uses":
		// sort.Sort(stacksByName(values))
		sort.Sort(stacksByCount(values))
	case "pages":
		// sort.Sort(stacksByName(values))
		sort.Sort(stacksBySourceCount(values))
	case "az":
		sort.Sort(stacksByName(values))
	}
	return values
}

func GetStackedEntries(game, level, show, search, sortBy, language string, user *User) []*StackedEntry {
	leveln, err := strconv.Atoi(level)
	if err != nil || leveln > 4 || leveln < 1 {
		leveln = 0
	}
	entries := GetEntriesAt(game, leveln, show, search, language, user)
	stacks := stackEntries(entries)
	return sortStacks(stacks, sortBy)
}

func (e *Entry) GetStackedEntry() *StackedEntry {
	entries := e.GetParts()
	stacked := stackEntries(entries)
	if len(stacked) == 0 {
		return nil
	}
	return stacked[0]
}

/* Stacked Translations */

func (se *StackedEntry) GetTranslations(language string) []*StackedTranslation {
	length := len(se.Entries)
	translations := make(map[string][]*Translation, 30)

	for _, entry := range se.Entries {
		entryTranslations := entry.GetTranslations(language)
		for _, translation := range entryTranslations {
			if _, ok := translations[translation.Translator]; !ok {
				translations[translation.Translator] = make([]*Translation, 0, length)
			}
			translations[translation.Translator] = append(translations[translation.Translator], translation)
		}
	}

	stackedTranslations := make([]*StackedTranslation, 0, len(translations))
	for _, parts := range translations {
		stacked := makeStackedTranslation(se, parts)
		if !stacked.Empty() {
			stackedTranslations = append(stackedTranslations, stacked)
		}
	}
	return stackedTranslations
}

func makeStackedTranslation(entry *StackedEntry, parts []*Translation) *StackedTranslation {
	isPreferred := false
	isConflicted := false
	language := parts[0].Language
	translator := parts[0].Translator
	for _, part := range parts {
		if part.IsPreferred {
			isPreferred = true
		}
		if part.IsConflicted {
			isConflicted = true
		}
	}
	
	text := make([]string, len(parts))
	for i, part := range parts {
		text[i] = part.Translation
	}
	fullText := strings.Join(text, "")

	stack := StackedTranslation{
		Entry:        entry,
		Language:     language,
		Translator:   translator,
		Parts:        parts,
		Count:        len(parts), // ???
		// SourceCount:  len(sources),
		FullText:     fullText,
		IsPreferred:  isPreferred,
		IsConflicted: isConflicted,
	}
	return &stack
}

type StackedTranslation struct {
	Entry        *StackedEntry
	Language     string
	Translator   string
	Parts        []*Translation
	Count        int
	// SourceCount  int
	FullText     string
	IsPreferred  bool
	IsConflicted bool
}

func (st *StackedTranslation) ID() uint64 {
	return hash64(st.Entry.FullText + " --- " + st.FullText)
}

func (st *StackedTranslation) Empty() bool {
	for _, part := range st.Parts {
		if part != nil && strings.TrimSpace(part.Translation) != "" {
			return false
		}
	}
	return true
}

func (se *StackedEntry) GetTranslationBy(language, translator string) *StackedTranslation {
	parts := make([]*Translation, len(se.Entries))
	for i, entry := range se.Entries {
		parts[i] = entry.GetTranslationBy(language, translator)
		if parts[i] == nil {
			parts[i] = &Translation{
				Entry:       *entry,
				Language:    language,
				Translation: "",
				Translator:  translator,
			}
		}
	}
	return makeStackedTranslation(se, parts)
}

func (se *StackedEntry) CountTranslations() map[string]int {
	entryCounts := make([]map[string]int, len(se.Entries))
	for i, entry := range se.Entries {
		entryCounts[i] = entry.CountTranslations()
	}

	langCounts := make(map[string]int, len(Languages))
	for _, lang := range Languages {
		min := 0
		for _, counts := range entryCounts {
			count := counts[lang]
			if count < min || min == 0 {
				min = count
			}
		}
		if min > 0 {
			langCounts[lang] = min
		}
	}
	return langCounts
}

// sort stacked entries by name
type stacksByName []*StackedEntry

func (this stacksByName) Len() int {
	return len(this)
}

func (this stacksByName) Less(i, j int) bool {
	collation := collate.New(language.English, collate.Loose)
	return collation.CompareString(this[i].FullText, this[j].FullText) > 0
	// return this[i].FullText < this[j].FullText
}

func (this stacksByName) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

// sort stacked entries by number of uses
type stacksByCount []*StackedEntry

func (this stacksByCount) Len() int {
	return len(this)
}

func (this stacksByCount) Less(i, j int) bool {
	return this[i].Count > this[j].Count
}

func (this stacksByCount) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}


// sort stacked entried by number of pages
type stacksBySourceCount []*StackedEntry

func (this stacksBySourceCount) Len() int {
	return len(this)
}

func (this stacksBySourceCount) Less(i, j int) bool {
	return this[i].SourceCount > this[j].SourceCount
}

func (this stacksBySourceCount) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}