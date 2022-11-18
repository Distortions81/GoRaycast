package main

import (
	"fmt"
	"image/color"
	"math"

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
	ebitenutil.DrawLine(screen,
		playerPhysics.Position.x*screenScale, playerPhysics.Position.y*screenScale,
		playerPhysics.Position.x*screenScale+playerPhysics.MovePos.x*32, playerPhysics.Position.y*screenScale+playerPhysics.MovePos.y*32,
		cYellow)
	ebitenutil.DrawCircle(screen, playerPhysics.Position.x*screenScale, playerPhysics.Position.y*screenScale, 8, cYellow)

	/* Draw rays */
	rayAngle := playerPhysics.Rotation
	var rayPos xycord
	var offset xycord
	for rayNum := 0; rayNum < 1; rayNum++ {
		/* Check Horizontal Lines */
		dof := 0
		arcTan := math.Atan(rayAngle)
		if rayAngle > onePi {
			rayPos.y = math.Floor(playerPhysics.Position.y/screenScale) * screenScale
			rayPos.x = (playerPhysics.Position.y-rayPos.y)*arcTan + playerPhysics.Position.y
			offset.y = -screenScale
			offset.x = -offset.y * arcTan
		} else if rayAngle < onePi {
			rayPos.y = math.Floor(playerPhysics.Position.y/screenScale)*screenScale + screenScale
			rayPos.x = (playerPhysics.Position.y-rayPos.y)*arcTan + playerPhysics.Position.y
			offset.y = screenScale
			offset.x = -offset.y * arcTan
		} else {
			rayPos.x = playerPhysics.Position.x
			rayPos.y = playerPhysics.Position.y
			dof = 8
		}
		for dof < 8 {
			if rayPos.x >= 0 && rayPos.x <= float64(mapSize.x) &&
				rayPos.y >= 0 && rayPos.y <= float64(mapSize.y) {
				red, _, _, _ := mapImg.At(int(rayPos.x), int(rayPos.y)).RGBA()
				if red > 0 {
					break /* hit wall */
				} else {
					/* next line */
					rayPos.x += offset.x
					rayPos.y += offset.y
					dof += 1
				}
			} else {
				break /* edge of map */
			}
		}
		ebitenutil.DrawLine(screen, playerPhysics.Position.x*screenScale, playerPhysics.Position.y*screenScale, rayPos.x*screenScale, rayPos.y*screenScale, cRed)

	}
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
