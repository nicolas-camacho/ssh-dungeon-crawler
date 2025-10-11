package game

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
	Type roomType
}

type model struct {
	state  gameState
	width  int
	height int

	loadingProgress float64

	menuCursor int

	playerMapX int
	playerMapY int
	worldMap   [][]*room
	stats      playerStats
}
