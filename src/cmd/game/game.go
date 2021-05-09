package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	SCREENWIDTH  = 1280
	SCREENHEIGHT = 720
	WALKPATH     = `Knight_02\02-Walk\`
)

var (
	images []*ebiten.Image
)

func init() {
	files, err := ioutil.ReadDir(WALKPATH)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file)
		f, err := os.Open(fmt.Sprint(WALKPATH, file.Name()))
		if err != nil {
			log.Fatal(err)
		}
		img, _, err := image.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
		image, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
		if err != nil {
			log.Fatal(err)
		}
		images = append(images, image)
	}
}

type Sprite struct {
	imageWidth  int
	imageHeight int
	x           int
	y           int
	vx          int
	vy          int
	frame       int
	frameCount  int
	hold        int
	flipped     bool
}

func (s *Sprite) Update() {
	s.hold = (s.hold + 1) % 3
	if s.hold == 0 {
		s.frame = (s.frame + 1) % s.frameCount
	}

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
		s.flipped = !s.flipped
	} else if my := SCREENHEIGHT - s.imageHeight; my <= s.y {
		s.y -= s.vy
		s.vy = 0
		s.vx = -3
		s.flipped = !s.flipped
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

	w, h := images[0].Size()
	x, y := rand.Intn(SCREENWIDTH-w), rand.Intn(SCREENHEIGHT-h)
	vx, vy := 3, 0
	g.sprite = &Sprite{
		imageWidth:  w / 4,
		imageHeight: h / 4,
		x:           x,
		y:           y,
		vx:          vx,
		vy:          vy,
		frameCount:  8,
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
	if g.sprite.flipped {
		g.op.GeoM.Scale(-.25, .25)
		g.op.GeoM.Translate(float64(g.sprite.imageWidth), 0)
	} else {
		g.op.GeoM.Scale(.25, .25)
	}

	g.op.GeoM.Translate(float64(g.sprite.x), float64(g.sprite.y))
	screen.DrawImage(images[g.sprite.frame], &g.op)
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
