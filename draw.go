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

	var op *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	op.GeoM.Scale(mapScale, mapScale)
	op.GeoM.Translate(float64(screenWidth-mapSize.x*mapScale), float64(screenHeight-mapSize.y*mapScale))
	op.Filter = ebiten.FilterNearest
	screen.DrawImage(mapImg, op)

	op.GeoM.Reset()
	op.GeoM.Translate(float64(screenWidth-mapSize.x*mapScale), float64(screenHeight-mapSize.y*mapScale))
	op.Filter = ebiten.FilterNearest
	screen.DrawImage(playerImg, op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
}
