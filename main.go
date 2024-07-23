package main

import (
	"flag"
	"log"
	"tower-defense/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	width := flag.Int("width", 80, "Board width in cells")
	height := flag.Int("height", 80, "Board height in cells")
	speed := flag.Int("speed", 15, "Ticks per second, min 0 max 60, + or - to adjust in game")
	debug := flag.Bool("debug", false, "Show debug info, D to toggle in game")
	lines := flag.Bool("lines", false, "Draw grid lines, L to toggle in game")

	flag.Parse()

	g, err := game.NewGame(*width, *height, *speed, *debug, *lines)
	if err != nil {
		log.Fatal(err)
	}
	err = g.Init()
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
