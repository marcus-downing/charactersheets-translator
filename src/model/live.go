package model

import (
	"encoding/json"
	"fmt"
)

type signal struct{}

func liveEntries() []*StackedEntry {
	strings := make(map[string]signal, 512)
	for _, game := range []string{"pathfinder", "dnd35"} {
		gameData := ReadGameData(game)
		if gameData != nil {
			// string += gameData skills | { ^ displayName }
			for _, skill := range gameData.Skills {
				strings[skill.SkillName()] = signal{}
			}

			// strings += gameData class | { ^ displayName }
			for _, class := range gameData.Classes {
				strings[class.Name+" Level"] = signal{}
			}
		}
	}

	entries := make([]*Entry, 0, len(strings))
	for str, _ := range strings {
		entry := Entry{str, str}
		entries = append(entries, &entry)
	}
	stacked := stackEntries(entries)
	fmt.Println("Found", len(entries), "entries to translate")
	return stacked
}

// GetLiveTranslations gives the translations needed by the Composer app
func GetLiveTranslations() []byte {
	entries := liveEntries()

	// translations := make([]*StackedTranslation, 0, len(entries)*len(Languages))

	var liveTranslations LiveTranslations
	liveTranslations.Languages = make([]LiveTranslationsLanguage, 0, len(Languages))
	for _, language := range Languages {
		languageTranslations := make([]LiveTranslationEntry, 0, len(entries))
		for _, entry := range entries {
			translations := entry.GetTranslations(language)
			selected := PickPreferredTranslation(entry.RankTranslations(translations, false))
			if selected != nil {
				for i, part := range entry.Entries {
					languageTranslations = append(languageTranslations, LiveTranslationEntry{
						Original:    part.Original,
						Translation: selected.Parts[i].Translation,
					})
				}
			}
		}
		// fmt.Println(" -", language, "-", len(languageTranslations), "translations")

		if len(languageTranslations) > 0 {
			liveTranslations.Languages = append(liveTranslations.Languages, LiveTranslationsLanguage{
				Name:         LanguagePaths[language],
				Translations: languageTranslations,
			})
		}
	}
	// fmt.Println("Exporting:", liveTranslations)
	return liveTranslations.export()
}

type LiveTranslations struct {
	Languages []LiveTranslationsLanguage `json:"languages"`
}

type LiveTranslationsLanguage struct {
	Name         string `json:"name"`
	Translations []LiveTranslationEntry `json:"translations"`
}

type LiveTranslationEntry struct {
	Original    string `json:"original"`
	Translation string `json:"translation"`
}

func (liveTranslations LiveTranslations) export() []byte {
	data, err := json.Marshal(liveTranslations)
	if err != nil {
		return nil
	}
	return data
}

func GetMasterInjectionEntries() []*StackedEntry {
	// entries := make()

	pathfinder := ReadGameData("pathfinder")
	// dnd35 := ReadGameData("dnd35")
/*
	pathfinderCharInfoPages := []string{
        "Pathfinder/Core/Character Info.ai", 
        "Pathfinder/Core/Animal Companion.ai",
        "Pathfinder/Core/Barbarian/Barbarian - Character Info.ai",
        "Pathfinder/Core/Ranger/Ranger - Character Info.ai",
        "Pathfinder/GM/NPC.ai",
        "Pathfinder/GM/NPC Group.ai",
        "Pathfinder/Archetypes/Druid/World Walker - Character Info.ai",
    }*/
	
	pathfinderSkills := make(map[string]signal, len(pathfinder.Skills))
	// string += gameData skills | { ^ displayName }
	for _, skill := range pathfinder.Skills {
		pathfinderSkills[skill.SkillName()] = signal{}
	}

	/*
	var pathfinderSkills = pathfinder skills ^ skillName
	*/

	return nil
}