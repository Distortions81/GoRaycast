package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Update()

	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterNearest}

	/*Draw walls */
	renderMap()
	screen.DrawImage(mapRender, op)

	/* Draw Player */
	op.GeoM.Reset()
	op.GeoM.Scale(2, 20)
	op.GeoM.Rotate(playerPhysics.Rotation)
	op.GeoM.Translate((playerPhysics.Position.x * screenScale), (playerPhysics.Position.y * screenScale))
	screen.DrawImage(lineImg, op)

	/* Draw rays */
	//

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
}

func renderMap() {
	if mapDirty {
		mapDirty = false
		mapRender.Fill(color.Transparent)

		op := &ebiten.DrawImageOptions{Filter: ebiten.FilterNearest}
		/* Draw walls */
		for x := 0; x < int(mapSize.x); x++ {
			for y := 0; y < int(mapSize.y); y++ {
				r, _, _, _ := mapImg.At(x, y).RGBA()
				if r > 0 {
					op.GeoM.Reset()
					op.GeoM.Scale(screenScale-1, screenScale-1)
					op.GeoM.Translate(float64(x*screenScale), float64(y*screenScale))
					mapRender.DrawImage(wallImg, op)
				}
			}
		}
	}
}
