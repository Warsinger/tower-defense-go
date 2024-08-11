package network

import (
	comp "tower-defense/components"

	"github.com/leap-fish/necs/esync"
)

func RegisterComponenets() {
	_ = esync.RegisterComponent(10, comp.TowerData{}, comp.Tower)
	_ = esync.RegisterComponent(11, comp.CreepData{}, comp.Creep)
	_ = esync.RegisterComponent(12, comp.PlayerData{}, comp.Player)
	_ = esync.RegisterComponent(13, comp.BulletData{}, comp.Bullet)
	_ = esync.RegisterComponent(14, comp.PositionData{}, comp.Position)
	_ = esync.RegisterComponent(15, comp.HealthData{}, comp.Health)
	_ = esync.RegisterComponent(16, comp.BoardData{}, comp.Board)
	_ = esync.RegisterComponent(17, comp.AttackData{}, comp.Attack)
	_ = esync.RegisterComponent(18, comp.SpriteRenderData{}, comp.SpriteRender)
	_ = esync.RegisterComponent(19, comp.PlayerRenderData{}, comp.PlayerRender)
	_ = esync.RegisterComponent(20, comp.RangeRenderData{}, comp.RangeRender)
	_ = esync.RegisterComponent(21, comp.InfoRenderData{}, comp.InfoRender)
	_ = esync.RegisterComponent(22, comp.BulletRenderData{}, comp.BulletRender)
	_ = esync.RegisterComponent(23, comp.LevelData{}, comp.Level)
	_ = esync.RegisterComponent(24, comp.BattleSceneState{}, comp.BattleState)
}

type ClientConnectMessage struct {
	Address string
}

type StartGameMessage struct{}

type CreepMessage struct {
	Count int
}
