package main

import (
	"fmt"
	"image"
	"log"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	SCREENWIDTH  = 1280
	SCREENHEIGHT = 720
)

var (
	ebitenImage *ebiten.Image
)

func init() {
	f, err := os.Open("robot.png")
	if err != nil {
		log.Fatal((err))
	}

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	oldebitenImage, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		log.Fatal((err))
	}
	ebitenImage = oldebitenImage
}

type Sprite struct {
	imageWidth  int
	imageHeight int
	x           int
	y           int
	vx          int
	vy          int
	angle       int
}

func (s *Sprite) Update() {
	s.x += s.vx
	s.y += s.vy
	if s.x < 0 {
		s.x -= s.vx
		s.vx = 0
		s.vy = -3
	} else if mx := SCREENWIDTH - s.imageWidth; mx <= s.x {
		s.x -= s.vx
		s.vx = 0
		s.vy = 3
	}
	if s.y < 0 {
		s.y -= s.vy
		s.vy = 0
		s.vx = 3
	} else if my := SCREENHEIGHT - s.imageHeight; my <= s.y {
		s.y -= s.vy
		s.vy = 0
		s.vx = -3
	}
}

type Game struct {
	sprite *Sprite
	op     ebiten.DrawImageOptions
	inited bool
}

func (g *Game) init() {
	defer func() {
		g.inited = true
	}()

	w, h := ebitenImage.Size()
	x, y := rand.Intn(SCREENWIDTH-w), rand.Intn(SCREENHEIGHT-h)
	vx, vy := 2, 0
	g.sprite = &Sprite{
		imageWidth:  w,
		imageHeight: h,
		x:           x,
		y:           y,
		vx:          vx,
		vy:          vy,
		angle:       0,
	}
}

func (g *Game) Update(img *ebiten.Image) error {
	if !g.inited {
		g.init()
	}
	g.sprite.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.op.GeoM.Reset()
	g.op.GeoM.Scale(1, 1)
	g.op.GeoM.Translate(float64(g.sprite.x), float64(g.sprite.y))
	screen.DrawImage(ebitenImage, &g.op)
	msg := fmt.Sprintf(`TPS: %0.2f
FPS: %0.2f
X: %v Y: %v`, ebiten.CurrentTPS(), ebiten.CurrentFPS(), g.sprite.x, g.sprite.y)
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return SCREENWIDTH, SCREENHEIGHT
}

func main() {
	ebiten.SetWindowSize(SCREENWIDTH, SCREENHEIGHT)
	ebiten.SetWindowTitle("Game")
	game := Game{}
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
