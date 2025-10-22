package game

import (
	"math/rand"
	"sort"

	"github.com/charmbracelet/bubbles/progress"
)

type combatActionState int

const (
	ActionSelect combatActionState = iota
	AttackSelect
	MagicSelect
	ItemSelect
	TargetSelect
)

type CombatEntity interface {
	GetName() string
	GetHP() int
	GetMaxHP() int
	GetSpeed() int
	TakeDamage(amount int)
	IsPlayer() bool
	ModifyStat(stat string, value int)
}

type Player struct {
	data        *playerData
	isDefending bool
	Attacks     []Attack
	Magics      []Magic
}

type Attack struct {
	Name    string   `json:"name"`
	Sides   int      `json:"sides"`
	Effects []Effect `json:"effects"`
}

type Magic struct {
	Name    string   `json:"name"`
	Sides   int      `json:"sides"`
	Cost    int      `json:"cost"`
	Effects []Effect `json:"effects"`
}

type Item struct {
	Name   string
	Effect string
	Value  int
}

type Effect struct {
	Target string `json:"target"`
	Stat   string `json:"stat"`
	Sides  int    `json:"sides"`
}

func (p *Player) GetName() string       { return "@YOU" }
func (p *Player) GetHP() int            { return p.data.stats.hp }
func (p *Player) GetMaxHP() int         { return 100 }
func (p *Player) GetSpeed() int         { return p.data.stats.speed }
func (p *Player) TakeDamage(amount int) { p.data.stats.hp -= amount }
func (p *Player) IsPlayer() bool        { return true }
func (p *Player) ModifyStat(stat string, value int) {
	switch stat {
	case "HP":
		p.data.stats.hp += value
		if p.data.stats.hp > p.GetMaxHP() {
			p.data.stats.hp = p.GetMaxHP()
		}
	case "Strength":
		p.data.stats.strength += value
	case "Defense":
		p.data.stats.defense += value
	case "Speed":
		p.data.stats.speed += value
	case "Magic":
		p.data.stats.magic += value
	}
}

type Foe struct {
	Name     string   `json:"name"`
	HP       int      `json:"hp"`
	MaxHP    int      `json:"maxHP"`
	Speed    int      `json:"speed"`
	Defense  int      `json:"defense"`
	Strength int      `json:"strength"`
	Attacks  []Attack `json:"attacks"`
}

func (e *Foe) GetName() string       { return e.Name }
func (e *Foe) GetHP() int            { return e.HP }
func (e *Foe) GetMaxHP() int         { return e.MaxHP }
func (e *Foe) GetSpeed() int         { return e.Speed }
func (e *Foe) TakeDamage(amount int) { e.HP -= amount }
func (e *Foe) IsPlayer() bool        { return false }
func (e *Foe) ModifyStat(stat string, value int) {
	switch stat {
	case "HP":
		e.HP += value
	case "Strength":
		e.Strength += value
	case "Defense":
		e.Defense += value
		if e.Defense < 0 {
			e.Defense = 0
		}
	case "Speed":
		e.Speed += value
	}
}

func newGoblin() *Foe {
	goblin := EnemyTemplates["goblin"]
	return &goblin
}

func calculateTurnOrder(player *Player, enemies []*Foe) []CombatEntity {
	entities := make([]CombatEntity, 0, len(enemies)+1)
	entities = append(entities, player)
	for _, e := range enemies {
		entities = append(entities, e)
	}

	sort.Slice(entities, func(i, j int) bool {
		return entities[i].GetSpeed() > entities[j].GetSpeed()
	})

	return entities
}

func newTestCombatState() *CombatState {
	playerData := &playerData{
		stats: playerStats{
			hp:       100,
			mana:     50,
			speed:    10,
			magic:    12,
			strength: 8,
			defense:  8,
		},
		inventory: map[string]int{
			"potion": 1,
		},
	}

	var playerAttacks []Attack
	for _, attack := range AttackTemplates {
		playerAttacks = append(playerAttacks, attack)
	}

	var playerMagics []Magic
	for _, magic := range MagicTemplates {
		playerMagics = append(playerMagics, magic)
	}

	playerEntity := &Player{
		data:    playerData,
		Attacks: playerAttacks,
		Magics:  playerMagics,
	}

	numEnemies := 1 + rand.Intn(3)
	enemies := make([]*Foe, numEnemies)
	for i := range enemies {
		enemies[i] = newGoblin()
	}

	turnOrder := calculateTurnOrder(playerEntity, enemies)

	enemyProgressBar := progress.New(
		progress.WithGradient(string(indigo), string(orange)),
		progress.WithoutPercentage(),
	)

	return &CombatState{
		player:                playerEntity,
		enemies:               enemies,
		turnOrder:             turnOrder,
		turnIndex:             0,
		actionState:           ActionSelect,
		isEnemyTurnInProgress: false,
		enemyActionProgress:   enemyProgressBar,
	}
}
