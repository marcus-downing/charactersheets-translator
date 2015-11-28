package model

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
	"path/filepath"
	"../config"
)

type GameData struct {
	Game string `json:"game"`
	Name string `json:"name"`

	Skills          []GameDataSkill `json:"skills"`
	CoreSkills      []string        `json:"coreSkills"`
	SummarySkills   []string        `json:"summarySkills"`
	KnowledgeSkills []string        `json:"knowledgeSkills"`
	AnimalSkills    []string        `json:"animalSkills"`

	Pages     []*GameDataPage `json:"pages"`
	GM        *GameGMData     `json:"gm"`
	Base      *GameDataClass  `json:"base"`
	Layout    [][]string
	Languages []*GameDataLanguage
	Books     []*struct {
		Name    string   `json:"name"`
		Classes []string `json:"classes"`
	}

	Classes []*GameDataClass `json:"classes"`
}

type GameDataSkill struct {
	Name               string `json:"name"`
	DisplayName        string `json:"displayName"`
	Ability            string `json:"ability"`
	UseUntrained       bool   `json:"useUntrained"`
	ArmourClassPenalty bool   `json:"acp"`
	NoRage             bool   `json:"noRage"`
	FavouredEnemy      bool   `json:"favouredEnemy"`
	FavouredTerrain    bool   `json:"favouredTerrain"`
	SubSkillOf         string `json:"subSkillOf"`
}

func (skill *GameDataSkill) SkillName() string {
	if skill.DisplayName == "" {
		return skill.Name
	}
	return skill.DisplayName
}

type GameDataPage struct {
	File     string `json:"file"`
	Slot     string `json:"slot"`
	Variant  string `json:"variant"`
	Name     string `json:"name"`
	Position int    `json:"position"`
}

type GameDataLanguage struct {
	Code      string `json:"code"`
	ShortCode string `json:"short"`
	Name      string `json:"name"`
	Ready     [4]int `json:"ready"`
}

type GameDataClass struct {
	Name     string           `json:"name"`
	Pages    []string         `json:"pages"`
	Skills   []string         `json:"skills"`
	Variants []*GameDataClass `json:"variants"`
}

type GameGMData struct {

}


func ReadGameData(game string) *GameData {
	filename := config.Config.PDF.Path + "/data/"+game+".json"
	path, err := filepath.Abs(filename)
	fmt.Println("Reading file:", path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file", err)
		return nil
	}

	var gameData GameData
	json.Unmarshal(data, &gameData)
	return &gameData
}
