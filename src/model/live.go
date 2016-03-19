package model

import (
	"encoding/json"
	"fmt"
)

type signal struct{}

var liveEntriesLit []string = []string{
	"Acrobatics",
	"Aegis Level",
	"Alchemist Level",
	"Appraise",
	"Arcanist Level",
	"Archivist Level",
	"Ardent Level",
	"Artificer Level",
	"Athletics",
	"Autohypnosis",
	"Balance",
	"Barbarian Level",
	"Bard Level",
	"Battle Dancer Level",
	"Beguiler Level",
	"Binder Level",
	"Bloodrager Level",
	"Bluff",
	"Brawler Level",
	"Cavalier Level",
	"Cleric Level",
	"Climb",
	"Concentration",
	"Control Shape",
	"Crusader Level",
	"Cryptic Level",
	"Death Master Level",
	"Decipher Script",
	"Diplomacy",
	"Disable Device",
	"Disarm Traps",
	"Disguise",
	"Divine Mind Level",
	"Dragon Shaman Level",
	"Dragonfire Adept Level",
	"Dread Level",
	"Dread Necromancer Level",
	"Druid Level",
	"Duskblade Level",
	"Eidolon Level",
	"Escape Artist",
	"Factotum Level",
	"Favoured Soul Level",
	"Fighter Level",
	"Finesse",
	"Fly",
	"Forgery",
	"Gather Information",
	"Gunslinger Level",
	"Handle Animal",
	"Heal",
	"Hexblade Level",
	"Hide",
	"High Guard Level",
	"Hunter Level",
	"Iaijutsu Focus",
	"Imperial Man-at-arms Level",
	"Incarnate Level",
	"Influence",
	"Inquisitor Level",
	"Intimidate",
	"Investigator Level",
	"Jester Level",
	"Jump",
	"Khalid Asad Level",
	"Knowledge (aeronautics)",
	"Knowledge (arcana)",
	"Knowledge (dungeoneering)",
	"Knowledge (engineering)",
	"Knowledge (geography)",
	"Knowledge (history)",
	"Knowledge (local)",
	"Knowledge (nature)",
	"Knowledge (nobility)",
	"Knowledge (planes)",
	"Knowledge (psionics)",
	"Knowledge (religion)",
	"Linguistics",
	"Listen",
	"Locate Traps",
	"Lurk Level",
	"Magus Level",
	"Marksman Level",
	"Martial Adept Level",
	"Martial Lore",
	"Monk Level",
	"Mountebank Level",
	"Move Silently",
	"Mythic Level",
	"Nature",
	"Ninja Level",
	"Open Lock",
	"Oracle Level",
	"Paladin Level",
	"Panther Warrior Level",
	"Perception",
	"Perform (act)",
	"Perform (comedy)",
	"Perform (dance)",
	"Perform (keyboard)",
	"Perform (oratory)",
	"Perform (percussion)",
	"Perform (string)",
	"Perform (wind)",
	"Perform (sing)",
	"Performance",
	"Prepared Level",
	"Prestige Level",
	"Priest Level",
	"Psicraft",
	"Psion Level",
	"Psychic Warrior Level",
	"Ranger Level",
	"Religion",
	"Ride",
	"Rogue Level",
	"Samurai Level",
	"Savant Level",
	"Scout Level",
	"Scry",
	"Search",
	"Sense Motive",
	"Sha'ir Level",
	"Shadowcaster Level",
	"Shaman Level",
	"Shugenja Level",
	"Skald Level",
	"Slayer Level",
	"Sleight of Hand",
	"Society",
	"Sorcerer Level",
	"Soulborn Level",
	"Soulknife Level",
	"Spellcaster Level",
	"Spellcraft",
	"Spellthief Level",
	"Spirit Shaman Level",
	"Spontaneous Level",
	"Spot",
	"Stealth",
	"Summoner Level",
	"Survival",
	"Swashbuckler Level",
	"Swim",
	"Swordsage Level",
	"Tactician Level",
	"Totemist Level",
	"Track",
	"Truenamer Level",
	"Truespeak",
	"Tumble",
	"Urban Druid Level",
	"Use Magic Device",
	"Use Rope",
	"Vitalist Level",
	"Warblade Level",
	"Warlock Level",
	"Warmage Level",
	"Warpriest Level",
	"Wild Empathy",
	"Wilder Level",
	"Witch Hunter Level",
	"Witch Level",
	"Wizard Level",
	"Wu Jen Level",
	"Wyrdcaster Level",
}

func liveEntries() []*StackedEntry {
	// scan for classes, skills etc
	strings := make(map[string]signal, 512)
	for _, game := range []string{"pathfinder", "dnd35"} {
		gameData := ReadGameData(game)
		if gameData != nil {
			// All skills
			// string += gameData skills | { ^ displayName }
			for _, skill := range gameData.Skills {
				strings[skill.SkillName()] = signal{}
			}

			// All classes
			// strings += gameData class | { ^ displayName }
			for _, class := range gameData.Classes {
				strings[class.Name+" Level"] = signal{}
			}
		}
	}
	// bring in the manual list as well, just to make sure
	for _, str := range liveEntriesLit {
		strings[str] = signal{}
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
		fmt.Println(" -", language, "-", len(languageTranslations), "translations")

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