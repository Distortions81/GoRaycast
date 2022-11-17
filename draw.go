package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Reset()
	op.GeoM.Scale(15, 2)
	op.GeoM.Rotate(playerRot)
	op.GeoM.Rotate(halfPi)
	op.GeoM.Translate((playerPos.x * screenScale), (playerPos.y * screenScale))
	screen.DrawImage(lineImg, op)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
}
