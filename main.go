package main

import (
	"fmt"
	
	"github.com/hajimehoshi/ebiten/v2"
	log "github.com/sirupsen/logrus"

	"github.com/vatine/popper2/pkg/game"
)

func main() {
	log.SetLevel(log.DebugLevel)
	g := game.NewGame()
	// g.TestSetup()

	fmt.Println(ebiten.RunGame(g))
}
