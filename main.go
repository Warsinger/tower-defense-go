package main

import (
	"flag"
	"log"
	"tower-defense/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	width := flag.Int("width", 600, "Board width in pixels")
	height := flag.Int("height", 800, "Board height in pixels")
	speed := flag.Int("speed", 60, "Ticks per second, min 0 max 60, + or - to adjust in game")
	debug := flag.Bool("debug", false, "Show debug info, D to toggle in game")
	towerLevel := flag.Int("level", 0, "Starting tower level to increase difficulty, 0 for default")
	computer := flag.Bool("computer", false, "Enable computer player")
	nosound := flag.Bool("nosound", false, "Turn off sound effects, S to toggle in game")

	flag.Parse()

	g, err := game.NewGame(*width, *height, *speed, *towerLevel, *debug, *computer, *nosound)
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
