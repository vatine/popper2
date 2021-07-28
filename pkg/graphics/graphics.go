package graphics

import (
	"image/color"
	"math"
	"math/rand"
	
	"github.com/hajimehoshi/ebiten/v2"
	log "github.com/sirupsen/logrus"
)

var (
	blank = color.RGBA{R: 0, G: 0, B: 0, A: 0}
)

// Basic game object, will essentially render as just a circle
type Sphere struct {
	x, y float64
	dx, dy float64
	r float64
	img *ebiten.Image
}

// Our main object-popping entity
type Explosion struct {
	x, y float64
	maxR, r float64
	innerR float64
	dr float64
}

// Render a sphere onto the game display
func (s *Sphere) Draw(t *ebiten.Image) {
	var g ebiten.GeoM

	g.Translate(s.x - s.r, s.y - s.r)
	opts := ebiten.DrawImageOptions{GeoM: g}
	t.DrawImage(s.img, &opts)
}

// Advance a sphere one step, bouncing against the borders if necessary.
func (s *Sphere) Step(w, h float64) {
	s.x += s.dx
	s.y += s.dy

	switch {
	case s.x - s.r < 0.0:
		s.dx = -s.dx
		s.x = s.r
	case s.x + s.r >= w:
		s.dx = -s.dx
		s.x = w - s.r
		
	}
	switch {
	case s.y - s.r < 0.0:
		fallthrough
	case s.y + s.r >= h:
		s.dy = -s.dy
		s.y += 2 * s.dy
	}

	if s.x < s.r || (s.x + s.r) > w {
		log.WithFields(log.Fields{
			"y": s.y,
			"x": s.x,
			"dx": s.dx,
			"dy": s.dy,
		}).Debug("x OOB")		
	}
	if s.y < s.r || (s.y + s.r) > h {
		log.WithFields(log.Fields{
			"y": s.y,
			"x": s.x,
			"dx": s.dx,
			"dy": s.dy,
		}).Debug("y OOB")		
	}
}

// Turn a sphere into an explosion.
func (s *Sphere) Explode() *Explosion {
	rv := new(Explosion)
	rv.x = s.x
	rv.y = s.y
	rv.maxR = s.r * 3.5
	rv.r = s.r
	rv.innerR = 0.0
	rv.dr = math.Sqrt(s.r)

	return rv
}

// Draw an explosion on the screen
func (e *Explosion) Draw(t *ebiten.Image) {
	ir := int(e.r + 0.5)
	bx := int(e.x)
	by := int(e.y)
	white := color.RGBA{R: 224, G: 224, B: 224, A: 128}
	for dx := -ir; dx <= ir; dx++ {
		for dy := -ir; dy <= ir; dy++ {
			fx := float64(bx + dx)
			fy := float64(by + dy)
			fdx := e.x - fx
			fdy := e.y - fy
			d2 := fdx*fdx + fdy*fdy

			if d2 <= (e.r * e.r) {
				if d2 >= (e.innerR * e.innerR) {
					t.Set(bx+dx, by+dy, white)
				}
			}
			
		}
	}
}

// Increase an explosion
func (e *Explosion) Step(w, h float64) {
	if e.r < e.maxR {
		e.r += e.dr
	}

	if e.r > e.maxR {
		e.r = e.maxR
	}

	if e.innerR > 0.0 {
		e.innerR += e.dr
	}

	if e.innerR == 0.0 && e.r > (0.5 * e.maxR) {
		e.innerR += e.dr		
	}
}

// Is the explosion done?
func (e *Explosion) Done() bool {
	return e.innerR >= e.maxR
}

// Clamp a float (v) between a desired minimum (l) and a maximum (h).
func clamp(v, l, h float64) float64 {
	if v >= h {
		return h
	}
	if v <= l {
		return l
	}
	return v
}

// Create a test sphere
func TestSphere(w, h float64, ix int) *Sphere {
	rv := new(Sphere)

	rv.r = 20
	if ix % 2 == 0 {
		rv.x = 20.0
		rv.dx = 1
	} else {
		rv.x = h - 20.0
		rv.x = -1
	}

	if ix & 2 == 0 {
		rv.y = 20
		rv.dy = 1.0
	} else {
		rv.y = h - 20
		rv.dy = -1
	}

	rv.paint()

	return rv
}

// Initial drawing of a sphere
func (s *Sphere) paint() {
	r := s.r
	red := clamp(rand.NormFloat64() * 15 + 175, 0, 255)
	green := clamp(rand.NormFloat64() * 15 + 175, 0, 255)
	blue := clamp(rand.NormFloat64() * 15 + 175, 0, 255)
	side := int(2*r+1)
	s.img = ebiten.NewImage(side, side)

	for x := 0; x < side; x++ {
		for y := 0; y < side; y++ {
			dx := r - float64(x)
			dy := r - float64(y)
			if (dx*dx)+(dy*dy) <= r*r {
				rf := ((dx*dx)+(dy*dy)) / (r*r)
				alpha := 192 + int(64 * rf)
				if alpha > 255 {
					alpha = 255
				}
				col := color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(alpha)}
				s.img.Set(x, y, col)
			}
		}
	}
}

// Return a multiplier and an offset for sphere radius.
func roundRadius(round int) (float64, float64) {
	shape := clamp(8.0 - (float64(round)/17), 0.01, 8)
	offset := clamp(30.0 - (float64(round)/8), 5, 30)

	return shape, offset
}

// Create a new sphere, in a random place, with random colours and a
// random direction of movement.
func NewSphere(w, h float64, round int) *Sphere {
	rv := new(Sphere)

	a := 2 * math.Pi * rand.Float64()
	s := clamp(rand.NormFloat64() * 1.5 + 5, 1.0, 9.0)
	shape, offset := roundRadius(round)
	r := clamp(rand.NormFloat64() * shape + offset, 7, 90.0)

	rv.x = (w - 2 * r) * rand.Float64() + r
	rv.y = (h - 2 * r) * rand.Float64() + r
	rv.dx = math.Cos(a) * s
	rv.dy = math.Sin(a) * s
	rv.r = r

	rv.paint()

	log.WithFields(log.Fields{
		"x": rv.x, "y": rv.y, "d": rv.dx, "dy": rv.dy, "r": rv.r,
	}).Debug("NewSphere")
	return rv
}

// Create a new explosion in a specific location. It will have a
// slightly randomised maximum radius.
func NewExplosion(x, y int) *Explosion {
	r := clamp(rand.NormFloat64() * 7 + 70, 50.0, 90.0)

	rv := new(Explosion)
	rv.r = 0.0
	rv.innerR = 0.0
	rv.dr = clamp(rand.NormFloat64() * 2 + 5, 3, 10)

	rv.maxR = r
	rv.x = float64(x)
	rv.y = float64(y)

	return rv
}

// Return true if a speger and explosion touch
func Intersect(s *Sphere, e *Explosion) bool {
	dx := s.x - e.x
	dy := s.y - e.y
	dr := s.r + e.r

	return (dx*dx)+(dy*dy) <= (dr*dr)
}
