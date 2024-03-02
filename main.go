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

func moveImage(positionChan chan imagePos, win windowDimensions) {

	var currentPos imagePos
	select {
	case currentPos = <-positionChan:
		// Successfully read the initial position
	default:
		// No initial position available, fallback to a default if necessary
		currentPos = imagePos{x: 0, y: 0} // Default position, adjust as needed
	}

	for {
		const progress = 0.5

		if currentPos.x >= float64(win.width) || currentPos.y >= float64(win.height) {
			// If the image has moved off the screen, reset the position
			currentPos.x = 0
			currentPos.y = 0
			// Once the image has fully moved across, you can stop the loop or reset progress
			positionChan <- imagePos{x: currentPos.x + (progress * 10), y: currentPos.y + (progress * 10)}
			// break

		}

		currentPos.x += (progress * 10)
		currentPos.y += (progress * 10)
		positionChan <- currentPos // Send the updated position

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

	positionChan := make(chan imagePos, 1)
	positionChan <- imagePos{x: 250, y: 150}

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
