package main

import (
	_ "image/png"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Declare your game struct
type Game struct {
	image       *ebiten.Image
	startTime   time.Time
	screenWidth float64
}

// Update proceeds the game state. Update is called every frame (1/60 [s]).
func (g *Game) Update() error {
	// Calculate elapsed time in seconds since the game started
	elapsed := time.Since(g.startTime).Seconds()

	// Assuming we want the image to move across the screen in 5 seconds
	// Calculate the new position based on the elapsed time
	xPos := (elapsed / 5) * g.screenWidth

	if xPos > g.screenWidth {
		// Restart the movement once the image reaches the end of the screen
		g.startTime = time.Now()
	}

	// Update the image position (this example keeps the Y position constant)
	g.image.Clear()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(xPos, 20) // Adjust Y position as needed
	g.image.DrawImage(g.image, op)

	return nil
}

// Draw draws the game screen. Draw is called every frame (1/60 [s]).
func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.image, nil)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	img, _, err := ebitenutil.NewImageFromFile("path/to/your/image.png")
	if err != nil {
		log.Fatalf("failed to load image: %v", err)
	}

	game := &Game{
		image:       img,
		startTime:   time.Now(),
		screenWidth: 800, // Assuming a screen width of 800 pixels
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Image Movement")

	// Start the game
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
