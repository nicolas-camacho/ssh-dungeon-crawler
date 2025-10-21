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
}

type Player struct {
	stats       *playerStats
	isDefending bool
	Attacks     []Attack
	Magics      []Magic
}

type Attack struct {
	Name  string
	Sides int
}

type Magic struct {
	Name  string
	Sides int
	Cost  int
}

func (p *Player) GetName() string       { return "@YOU" }
func (p *Player) GetHP() int            { return p.stats.hp }
func (p *Player) GetMaxHP() int         { return 100 }
func (p *Player) GetSpeed() int         { return p.stats.speed }
func (p *Player) TakeDamage(amount int) { p.stats.hp -= amount }
func (p *Player) IsPlayer() bool        { return true }

type Foe struct {
	Name    string
	HP      int
	MaxHP   int
	Speed   int
	Attack  int
	Defense int
}

func (e *Foe) GetName() string       { return e.Name }
func (e *Foe) GetHP() int            { return e.HP }
func (e *Foe) GetMaxHP() int         { return e.MaxHP }
func (e *Foe) GetSpeed() int         { return e.Speed }
func (e *Foe) TakeDamage(amount int) { e.HP -= amount }
func (e *Foe) IsPlayer() bool        { return false }

func newGoblin() *Foe {
	return &Foe{
		Name:    "Goblin",
		HP:      20,
		MaxHP:   20,
		Speed:   8,
		Attack:  5,
		Defense: 2,
	}
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
	playerStats := &playerStats{
		hp:       100,
		mana:     50,
		speed:    10,
		magic:    12,
		strength: 8,
		defense:  8,
	}
	playerEntity := &Player{
		stats: playerStats,
		Attacks: []Attack{
			{Name: "Slash", Sides: 4},
			{Name: "Final Slash", Sides: 6},
		},
		Magics: []Magic{
			{Name: "Fireball", Sides: 5, Cost: 10},
			{Name: "Firestorm", Sides: 8, Cost: 25},
		},
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
