package components

import (
	"tower-defense/assets"
	"tower-defense/config"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
)

type BattleSceneState struct {
	GameOver bool
	Paused   bool
}

var BattleState = donburi.NewComponentType[BattleSceneState]()

func (bss *BattleSceneState) Draw(screen *ebiten.Image, width, height float64, config *config.ConfigData, gameStats *GameStats) {
	var nextY float64
	if bss.GameOver {
		str := "GAME OVER"
		nextY = DrawTextLines(screen, assets.ScoreFace, str, width, height/2, text.AlignCenter, text.AlignCenter)

		str = "Press R to reset game"
		_ = DrawTextLines(screen, assets.InfoFace, str, width, nextY, text.AlignCenter, text.AlignStart)

	} else if bss.Paused {
		str := "PAUSED"
		_ = DrawTextLines(screen, assets.ScoreFace, str, width, height/2, text.AlignCenter, text.AlignCenter)
	}

	if gameStats != nil && config.ShowStats {
		str := gameStats.StatsLines(" ", !bss.GameOver, true, "High", "Player")
		_ = DrawTextLines(screen, assets.InfoFace, str, width, 450, text.AlignStart, text.AlignStart)
	}
}
