package main

import (
	"flag"
	"log"
	"tower-defense/config"
	"tower-defense/game"
	"tower-defense/strategy"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	width := flag.Int("width", 600, "Board width in pixels")
	height := flag.Int("height", 800, "Board height in pixels")
	speed := flag.Int("speed", 60, "Ticks per second, min 0 max 60, + or - to adjust in game")
	debug := flag.Bool("debug", false, "Show debug info, D to toggle in game")
	towerLevel := flag.Int("level", 0, "Starting tower level to increase difficulty, 0 for default")
	computer := flag.Bool("computer", false, "Enable computer player")
	computerLevel := flag.Int("complevel", 3, "Computer player difficulty level [1 slowest, 2 slow, 3 normal, 4 fast, 5 fastest]")
	nosound := flag.Bool("nosound", false, "Turn off sound effects, S to toggle in game")
	balancePath := flag.String("balance", "", "Path to game balance JSON config, empty for default")

	flag.Parse()

	strategy.SetComputerLevel(*computerLevel)
	balance, err := config.LoadBalance(*balancePath)
	if err != nil {
		log.Fatal(err)
	}
	g, err := game.NewGame(*width, *height, *speed, *towerLevel, *debug, *computer, *nosound, balance)
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
