package strategy

import (
	"fmt"
	"math"
	comp "tower-defense/components"
	"tower-defense/config"
	"tower-defense/util"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

const (
	towersPerRow   = 7
	towerWidth     = 48
	towerHeight    = 48
	halfTowerWidth = towerWidth / 2
	laneSpacing    = 41
	printTries     = false
	playSound      = false
)

// laneSpacing pixels between towers, towersPerRow towers across starting
var (
	lanes      = makeLanes()
	TimeScaler = 60
)

func SetComputerLevel(computerLevel int) {
	if computerLevel <= 1 {
		TimeScaler = 90
	} else if computerLevel == 2 {
		TimeScaler = 75
	} else if computerLevel == 3 {
		TimeScaler = 60
	} else if computerLevel == 4 {
		TimeScaler = 45
	} else if computerLevel == 5 {
		TimeScaler = 30
	} else if computerLevel >= 5 {
		TimeScaler = 1
	}
}

func Update(world donburi.World) (bool, error) {
	// can perform only one action per tick, even this might be too fast so maybe move this into the game speed updates
	pe := comp.Player.MustFirst(world)
	player := comp.Player.Get(pe)
	be := comp.Board.MustFirst(world)
	board := comp.Board.Get(be)
	debug := config.GetConfig(world).Debug

	query := donburi.NewQuery(
		filter.Or(
			filter.Contains(comp.Creep),
			filter.Contains(comp.Tower),
		),
	)

	count := query.Count(world)
	towers := make([]*donburi.Entry, 0, count)
	creeps := make([]*donburi.Entry, 0, count)

	query.Each(world, func(entry *donburi.Entry) {
		if entry.HasComponent(comp.Tower) {
			towers = append(towers, entry)
		} else if entry.HasComponent(comp.Creep) {
			creeps = append(creeps, entry)
		}

	})

	// if we have < the number of towers per row then check for a creep coming down and put a tower below it
	for _, creepEntry := range creeps {
		pt := util.MidpointRect(comp.GetRect(creepEntry))
		lane := findLane(lanes, pt.X)
		if lane != -1 {
			placed, err := player.TryPlaceTower(world, lane, board.Height/2, playSound, printTries)
			if err != nil {
				return false, err
			}
			if placed {
				if debug {
					fmt.Printf("Placed tower below creep at %v, %v\n", pt, lane)
				}
				return true, nil
			}
		}
	}

	// while we have fewer than N towers, don't you dare upgrade
	allowUpgrades := len(towers) >= towersPerRow && player.Money >= 75

	// if we have towers, if any need healing badly then heal them if < N or upgrade if >=N (and we have enough money)
	lowestHealthTower := findLowestHealthTower(towers)
	lowestLevelTower := findLowestLevelTower(towers)

	if lowestHealthTower != nil {
		// having found the lowest health and lowest level tower, they are the same the just upgrade, unless we don't have enough money
		if lowestHealthTower.Entity() == lowestLevelTower.Entity() {
			if allowUpgrades && player.TryUpgradeTower(lowestHealthTower, playSound, printTries) {
				if debug {
					fmt.Printf("Upgraded lowest health/level tower\n")
				}
				return true, nil
			}
			if player.TryHealTower(lowestHealthTower, playSound, printTries) {
				if debug {
					fmt.Printf("Healed lowest health tower\n")
				}
				return true, nil
			}
		}
		//if not the same then if the level of the lowest level tower is 2 less than the lowest health tower (or lowest health is max level)
		// then just heal the lowest health, otherwise upgrade it
		levelLevel := comp.Level.Get(lowestLevelTower)
		levelHealth := comp.Level.Get(lowestHealthTower)
		if levelLevel.Level >= levelHealth.Level+2 || levelHealth.Level == comp.GetMaxTowerLevel(world) {
			if player.TryHealTower(lowestHealthTower, playSound, printTries) {
				if debug {
					fmt.Printf("Healed lowest level tower\n")
				}
				return true, nil
			}
		} else {
			if allowUpgrades && player.TryUpgradeTower(lowestHealthTower, playSound, printTries) {
				if debug {
					fmt.Printf("Upgraded lowest health tower\n")
				}
				return true, nil
			}
		}
	}

	// if we have enough towers, upgrade the lowest tower (if we have enough money)
	if allowUpgrades {
		lowestLevelTower := findLowestLevelTower(towers)

		if player.TryUpgradeTower(lowestLevelTower, playSound, printTries) {
			if debug {
				fmt.Printf("Upgraded lowest health tower\n")
			}
			return true, nil
		}
	} else if lowestHealthTower != nil {
		if player.TryHealTower(lowestHealthTower, playSound, printTries) {
			if debug {
				fmt.Printf("Healed lowest health tower 2\n")
			}
			return true, nil
		}
	}

	// later game if we are full on towers and full on levels then start additional rows (up to 4) of towers to upgrade
	if player.Money > 150 && len(towers) >= towersPerRow {
		for i := 1; i <= 3; i++ {
			newY := board.Height/2 + i*(towerHeight+10)
			for _, lane := range lanes {
				placed, err := player.TryPlaceTower(world, lane, newY, playSound, printTries)
				if err != nil {
					return false, err
				}
				if placed {
					if debug {
						fmt.Printf("Placed tower on row %d at %v, %v\n", i, lane, newY)
					}
					return true, nil
				}
			}
		}
	}

	// TODO if multiplayer consider sending a creep over

	if debug {
		fmt.Printf("Nothing on this tick\n")
	}
	return false, nil
}

func findLowestHealthTower(towers []*donburi.Entry) *donburi.Entry {
	var lowestHealthTower *donburi.Entry
	var lowestHealth int = math.MaxInt
	for _, towerEntry := range towers {
		health := comp.Health.Get(towerEntry)
		percentHealth := float32(health.Health) / float32(health.MaxHealth)
		if percentHealth < 0.25 && health.Health <= lowestHealth {
			// find the tower with the lowest health below 25%
			lowestHealthTower = towerEntry
			lowestHealth = health.Health
		}
	}
	return lowestHealthTower
}

func findLowestLevelTower(towers []*donburi.Entry) *donburi.Entry {
	var lowestLevelTower *donburi.Entry
	var lowestLevel int = math.MaxInt
	for _, towerEntry := range towers {
		level := comp.Level.Get(towerEntry)
		if level.Level < lowestLevel {
			lowestLevelTower = towerEntry
			lowestLevel = level.Level
		}
	}
	return lowestLevelTower
}

func makeLanes() []int {
	lanes := make([]int, towersPerRow)
	lanes[0] = halfTowerWidth + 10

	for i := 1; i < len(lanes); i++ {
		lanes[i] = lanes[i-1] + towerWidth + laneSpacing
	}
	var sign int

	// order lanes starting from the middle position
	orderedlanes := make([]int, towersPerRow)
	mid := len(orderedlanes) / 2
	var offset = 0
	if len(lanes)%2 == 0 {
		offset = 1
	}
	for i := 0; i < len(orderedlanes); i++ {
		if i%2 == 0 {
			sign = -1
		} else {
			sign = 1
		}
		mid = mid + i*sign
		orderedlanes[i] = lanes[mid-offset]
	}
	return orderedlanes
}

// find the lane that is within towerWidth of the creep
func findLane(lanes []int, x int) int {
	for _, lane := range lanes {
		if util.Abs(x-lane) < towerWidth {
			return lane
		}
	}
	return -1
}
