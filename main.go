package main

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type imagePos struct {
	x float64
	y float64
}

type windowDimensions struct {
	width  int
	height int
}

type Game struct {
	win          windowDimensions
	img          *ebiten.Image
	positionChan chan imagePos
	imgPosition  imagePos
}

func (g *Game) Update() error {
	select {
	case newPos := <-g.positionChan:
		// Update position from channel
		g.imgPosition.x = newPos.x
		g.imgPosition.y = newPos.y
	default:
		// If no new position, do nothing
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// Use the updated position to draw the image
	op.GeoM.Translate(g.imgPosition.x, g.imgPosition.y)
	screen.DrawImage(g.img, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.win.width, g.win.height
}

func moveImage(positionChan chan imagePos, scr windowDimensions) {
	startTime := time.Now()
	totalDuration := 10 * time.Second // Total time to move across the screen

	for {
		elapsed := time.Since(startTime)
		progress := elapsed.Seconds() / totalDuration.Seconds()

		if progress >= 1.0 {
			// Once the image has moved across, you can stop the loop or reset progress
			progress = 1.0
			positionChan <- imagePos{x: float64(scr.width) * progress, y: float64(scr.height) * progress}
			break // Stop after one move; remove this if continuous movement is desired
		}

		positionChan <- imagePos{x: float64(scr.width) * progress, y: float64(scr.height) * progress}
		time.Sleep(16 * time.Millisecond) // About 60 FPS
	}
}

func main() {

	win := windowDimensions{width: 1024, height: 768}

	ebiten.SetWindowSize(win.width, win.height)
	ebiten.SetWindowTitle("Move Image")

	img, _, err := ebitenutil.NewImageFromFile("C:\\Users\\simonlomax\\Pictures\\sprites\\truck.png")
	if err != nil {
		log.Fatalf("failed to load image: %v", err)
	}

	positionChan := make(chan imagePos)

	g := &Game{
		win:          win,
		img:          img,
		positionChan: positionChan,
	}

	go moveImage(positionChan, win)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
