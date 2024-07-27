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
	server := flag.Int("server", 53838, "Startup multiplayer server on port")
	client := flag.String("client", "localhost:53838", "For multiplayer client connect to server:port")

	flag.Parse()

	g, err := game.NewGame(*width, *height, *speed, *debug)
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
