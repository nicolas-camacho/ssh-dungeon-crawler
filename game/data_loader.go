package game

import (
	"encoding/json"
	"os"
)

var (
	EnemyTemplates  map[string]Foe
	AttackTemplates map[string]Attack
	MagicTemplates  map[string]Magic
	ItemTemplates   map[string]Item
)

func LoadGameData() error {
	if err := loadFile("data/enemies.json", &EnemyTemplates); err != nil {
		return err
	}
	if err := loadFile("data/attacks.json", &AttackTemplates); err != nil {
		return err
	}
	if err := loadFile("data/magics.json", &MagicTemplates); err != nil {
		return err
	}
	if err := loadFile("data/items.json", &ItemTemplates); err != nil {
		return err
	}

	return nil
}

func loadFile(path string, target any) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}
