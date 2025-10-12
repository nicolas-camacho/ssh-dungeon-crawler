package game

import (
	"github.com/charmbracelet/bubbles/progress"
)

type gameState int

const (
	stateLoading gameState = iota
	stateMenu
	stateGame
)

const playerArt = `👤YOU`

type playerStats struct {
	hp, mana, speed, magic, strength int
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

type model struct {
	state  gameState
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
}
