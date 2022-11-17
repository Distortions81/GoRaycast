package main

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Update()

	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Reset()
	op.GeoM.Scale(2, 20)
	op.GeoM.Rotate(playerPhysics.Rotation)
	op.GeoM.Translate((playerPhysics.Position.x * screenScale), (playerPhysics.Position.y * screenScale))
	screen.DrawImage(lineImg, op)

	/* Draw rays */
	var rx, ry float64
	aTan := math.Atan(playerPhysics.Rotation)

	/* Looking down */
	if playerPhysics.Rotation > onePi {

		ry = playerPhysics.Position.y
		rx = (playerPhysics.Position.y-ry)*aTan + playerPhysics.Position.x
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
}
