package main

import (
	"flag"
	"fmt"
	
	"github.com/hajimehoshi/ebiten/v2"
	log "github.com/sirupsen/logrus"

	"github.com/vatine/popper2/pkg/game"
)

func main() {
	verbose := flag.Bool("v", false, "Verbose (dbug) logging")
	testMode := flag.Bool("d", false, "Debug (test) mode")
	g := game.NewGame()
	if *verbose {
		log.SetLevel(log.DebugLevel)
	}
	if *testMode {
		g.TestSetup()
	}

	fmt.Println(ebiten.RunGame(g))
}
