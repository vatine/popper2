package main

import (
	"flags"
	"fmt"
	
	"github.com/hajimehoshi/ebiten/v2"
	log "github.com/sirupsen/logrus"

	"github.com/vatine/popper2/pkg/game"
)

func main() {
	verbose := flags.Bool("v", false, "Verbose (dbug) logging")
	testMode := flags.Bool("d", false, "Debug (test) mode")
	g := game.NewGame()
	if *verbose {
		log.SetLevel(log.DebugLevel)
	}
	if *testMode {
		g.TestSetup()
	}

	fmt.Println(ebiten.RunGame(g))
}
