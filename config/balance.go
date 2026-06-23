package config

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"

	"github.com/yohamta/donburi"
)

//go:embed default_balance.json
var balanceFS embed.FS

type BalanceData struct {
	Player      PlayerBalance      `json:"player"`
	Tower       TowerBalance       `json:"tower"`
	Creep       CreepBalance       `json:"creep"`
	SuperCreep  SuperCreepBalance  `json:"superCreep"`
	Wave        WaveBalance        `json:"wave"`
	Multiplayer MultiplayerBalance `json:"multiplayer"`
}

type PlayerBalance struct {
	StartingMoney          int `json:"startingMoney"`
	Health                 int `json:"health"`
	AttackPower            int `json:"attackPower"`
	AttackRange            int `json:"attackRange"`
	AttackCooldown         int `json:"attackCooldown"`
	CreepLevelTowerLevels  int `json:"creepLevelTowerLevels"`
	MaxTowerInitialLevel   int `json:"maxTowerInitialLevel"`
	MaxTowerLevelsPerBonus int `json:"maxTowerLevelsPerBonus"`
}

type TowerBalance struct {
	DefaultType              string         `json:"defaultType"`
	Costs                    map[string]int `json:"costs"`
	HealCostDivisor          int            `json:"healCostDivisor"`
	Health                   int            `json:"health"`
	AttackPower              int            `json:"attackPower"`
	AttackRange              int            `json:"attackRange"`
	AttackCooldown           int            `json:"attackCooldown"`
	InitialLevel             int            `json:"initialLevel"`
	UpgradeMaxHealthAdd      int            `json:"upgradeMaxHealthAdd"`
	UpgradePowerLevelDivisor int            `json:"upgradePowerLevelDivisor"`
	UpgradeRangeAdd          int            `json:"upgradeRangeAdd"`
	UpgradeCooldownReduction int            `json:"upgradeCooldownReduction"`
	UpgradeMinCooldown       int            `json:"upgradeMinCooldown"`
}

type CreepBalance struct {
	BigCreepChance                  float32 `json:"bigCreepChance"`
	BigCreepChoice                  int     `json:"bigCreepChoice"`
	BigCreepAugment                 int     `json:"bigCreepAugment"`
	SmallCreepFirstChoice           int     `json:"smallCreepFirstChoice"`
	SmallCreepVariants              int     `json:"smallCreepVariants"`
	BaseVelocityY                   int     `json:"baseVelocityY"`
	HealthBase                      int     `json:"healthBase"`
	HealthAugmentMultiplier         int     `json:"healthAugmentMultiplier"`
	HealthLevelDivisor              int     `json:"healthLevelDivisor"`
	AttackPowerBase                 int     `json:"attackPowerBase"`
	AttackPowerLevelOffset          int     `json:"attackPowerLevelOffset"`
	AttackPowerLevelDivisor         int     `json:"attackPowerLevelDivisor"`
	AttackRangeBase                 int     `json:"attackRangeBase"`
	AttackRangeAugmentMultiplier    int     `json:"attackRangeAugmentMultiplier"`
	AttackCooldownBase              int     `json:"attackCooldownBase"`
	AttackCooldownAugmentMultiplier int     `json:"attackCooldownAugmentMultiplier"`
	ScoreValueBase                  int     `json:"scoreValueBase"`
}

type SuperCreepBalance struct {
	VelocityX      int `json:"velocityX"`
	VelocityY      int `json:"velocityY"`
	ScoreValue     int `json:"scoreValue"`
	Health         int `json:"health"`
	AttackPower    int `json:"attackPower"`
	AttackRange    int `json:"attackRange"`
	AttackCooldown int `json:"attackCooldown"`
}

type WaveBalance struct {
	SpawnBorder          int                `json:"spawnBorder"`
	MaxCreepTimer        int                `json:"maxCreepTimer"`
	StartCreepTimer      int                `json:"startCreepTimer"`
	MinCreepTick         int                `json:"minCreepTick"`
	MaxCreepCount        int                `json:"maxCreepCount"`
	TimerLevelDivisor    int                `json:"timerLevelDivisor"`
	ExtraCreepLevelWaves int                `json:"extraCreepLevelWaves"`
	LevelBumpDivisor     int                `json:"levelBumpDivisor"`
	SpawnIncomePerCreep  int                `json:"spawnIncomePerCreep"`
	OverflowIncome       int                `json:"overflowIncome"`
	SpawnChances         []CreepSpawnChance `json:"spawnChances"`
}

type CreepSpawnChance struct {
	Count  int     `json:"count"`
	Chance float32 `json:"chance"`
}

type MultiplayerBalance struct {
	SuperCreepCost     int `json:"superCreepCost"`
	SuperCreepCooldown int `json:"superCreepCooldown"`
}

var Balance = donburi.NewComponentType[BalanceData]()
var defaultBalance = mustLoadDefaultBalance()

func LoadBalance(path string) (*BalanceData, error) {
	if path == "" {
		return parseBalanceFile("embedded default balance", mustReadDefaultBalance())
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parseBalanceFile(path, bytes)
}

func DefaultBalance() *BalanceData {
	return defaultBalance
}

func NewBalance(world donburi.World, balance *BalanceData) *BalanceData {
	if balance == nil {
		balance = DefaultBalance()
	}
	entity := world.Create(Balance)
	entry := world.Entry(entity)

	Balance.Set(entry, balance)
	return Balance.Get(entry)
}

func GetBalance(world donburi.World) *BalanceData {
	entry, ok := Balance.First(world)
	if !ok {
		return DefaultBalance()
	}
	return Balance.Get(entry)
}

func mustLoadDefaultBalance() *BalanceData {
	balance, err := parseBalanceFile("embedded default balance", mustReadDefaultBalance())
	if err != nil {
		panic(err)
	}
	return balance
}

func mustReadDefaultBalance() []byte {
	bytes, err := balanceFS.ReadFile("default_balance.json")
	if err != nil {
		panic(err)
	}
	return bytes
}

func parseBalanceFile(name string, bytes []byte) (*BalanceData, error) {
	var balance BalanceData
	if err := json.Unmarshal(bytes, &balance); err != nil {
		return nil, fmt.Errorf("parse %s: %w", name, err)
	}
	return &balance, nil
}
