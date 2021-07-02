package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	log "github.com/sirupsen/logrus"

	"github.com/vatine/popper2/pkg/graphics"
)

var black = color.RGBA{R: 0, G: 0, B: 0, A: 255}
var bomb *ebiten.Image

// Main game container
type Game struct {
	round      int
	spheres    []*graphics.Sphere
	explosions []*graphics.Explosion
	available  int
	left       int
	w, h       int
}

// Create a new "designed for play" game.
func NewGame() *Game {
	g := new(Game)

	g.w = 800
	g.h = 600
	g.available = 1

	g.NewRound()

	if bomb == nil {
		bomb = ebiten.NewImage(10, 10)
		red := color.RGBA{R: 255, A: 255}
		for x := 0; x < 10; x++ {
			for y := 0; y < 10; y++ {
				dx := (4 -x)
				if dx < 0 {
					dx = -dx
				}
				dy := (4 - y)
				if dy < 0 {
					dy = -dy
				}
				if dx+dy < 5 {
					bomb.Set(x, y, red)
				}
			}
		}
	}
	
	return g
}

// Convert a game to a dedicated test setup.
func (g *Game) TestSetup() {
	g.spheres = []*graphics.Sphere{
		graphics.TestSphere(320, 240, 0),
		graphics.TestSphere(320, 240, 1),
		graphics.TestSphere(320, 240, 2),
		graphics.TestSphere(320, 240, 3),
	}
}

// Compute fib(n)
func fib(n int) int {
	a, b := 1, 1
	for i := 0; i < n; i++ {
		a, b = b, a+b
	}
	return b
}

// Create a new round. This basically increases the number of
// sopheres, and occasionally increases the number of explosions
// available.
func (g *Game) NewRound() {
	log.WithFields(log.Fields{
		"round":      g.round,
		"spheres":    len(g.spheres),
		"explosions": len(g.explosions),
		"available": g.available,
	}).Debug("NewRound start")
	g.round++
	g.spheres = []*graphics.Sphere{}
	spheres := g.round
	for i := 0; i < spheres; i++ {
		g.spheres = append(g.spheres, graphics.NewSphere(float64(g.w), float64(g.h)))
	}
	g.explosions = []*graphics.Explosion{}
	if fib(g.available) <= g.round {
		g.available++
	}
	g.left = g.available
	log.WithFields(log.Fields{
		"round":      g.round,
		"spheres":    len(g.spheres),
		"explosions": len(g.explosions),
		"available": g.available,
	}).Debug("Newround end")
}

// Draw the game.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(black)

	for _, s := range g.spheres {
		s.Draw(screen)
	}

	for _, e := range g.explosions {
		e.Draw(screen)
	}

	for l := 0; l < g.left; l++ {
		var g ebiten.GeoM
		
		g.Translate(float64(12*l), 5)
		opts := ebiten.DrawImageOptions{GeoM: g}
		screen.DrawImage(bomb, &opts)
	}
}

// Is the game over? This means:
//  * There are spheres bouncing around
//  * There are no active explosions
//  * There are no bombs/explosions available
func (g *Game) done() bool {
	if len(g.explosions) > 0 {
		return false
	}
	if g.left > 0 {
		return false
	}
	if len(g.spheres) > 0 {
		return true
	}

	return false
}

// Is it time for the next round?
// This basically means:
//  * No active explosions
//  * No remaining spheres
func (g *Game) roundOver() bool {
	if len(g.spheres) != 0 {
		return false
	}
	if len(g.explosions) != 0 {
		return false
	}

	return true
}

// Compute the next frame
func (g *Game) Update() error {
	if g.done() {
		return fmt.Errorf("Game over, man")
	}

	fw := float64(g.w)
	fh := float64(g.h)
	for _, s := range g.spheres {
		s.Step(fw, fh)
	}
	var newEx []*graphics.Explosion
	for _, e := range g.explosions {
		if !e.Done() {
			e.Step(fw, fh)
			newEx = append(newEx, e)
		}
	}
	g.explosions = newEx

	var newSpheres []*graphics.Sphere
	for _, s := range g.spheres {
		explodes := false
		for _, e := range g.explosions {
			if graphics.Intersect(s, e) {
				explodes = true
			}
		}
		if explodes {
			g.explosions = append(g.explosions, s.Explode())
		} else {
			newSpheres = append(newSpheres, s)
		}
	}
	g.spheres = newSpheres

	if g.roundOver() {
		g.NewRound()
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.explosions = append(g.explosions, graphics.NewExplosion(x, y))
		g.left--
	}

	return nil
}

// Return a static layout
func (g *Game) Layout(a, b int) (int, int) {
	return g.w, g.h
}
