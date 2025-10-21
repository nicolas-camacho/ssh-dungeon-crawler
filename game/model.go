package game

import (
	"github.com/charmbracelet/bubbles/progress"
)

type GameState int

const (
	StateLoading GameState = iota
	StateMenu
	StateGame
	StateCombat
)

const playerArt = `👤YOU`

type playerStats struct {
	hp, mana, speed, magic, strength, defense int
}

type roomType int

const (
	Empty roomType = iota
	Enemy
	Tresure
	Shop
	StairsUp
	StairsDown
)

type room struct {
	Type    roomType
	Visited bool
}

type floor struct {
	worldMap [][]*room
}

type CombatState struct {
	player                *Player
	enemies               []*Foe
	turnOrder             []CombatEntity
	turnIndex             int
	actionState           combatActionState
	actionCursor          int
	subActionCursor       int
	targetCursor          int
	isEnemyTurnInProgress bool
	enemyActionProgress   progress.Model
}

type model struct {
	state  GameState
	width  int
	height int
	styles styles

	progress progress.Model

	menuCursor int

	floors       []floor
	currentFloor int
	playerMapX   int
	playerMapY   int
	stats        playerStats

	combat *CombatState
}
