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

	/* Draw rays */
	rayAngle := playerPhysics.Rotation
	var rayPos xycord
	var offset xycord
	quant := 1.0
	for rayNum := 0; rayNum < 30; rayNum++ {
		/* Check Horizontal Lines */
		dof := 0
		rayAngle = (rayAngle - 0.523599) + (float64(rayNum) * 0.0349066)
		if rayAngle > twoPi {
			rayAngle -= twoPi
		} else if rayAngle < 0 {
			rayAngle += twoPi
		}

		aTan := -1 / math.Tan(rayAngle)
		if rayAngle > onePi {
			rayPos.y = playerPhysics.Position.y + quant
			rayPos.x = (playerPhysics.Position.y-rayPos.y)*aTan + playerPhysics.Position.x
			offset.y = -quant
			offset.x = -offset.y * aTan
		} else if rayAngle < onePi {
			rayPos.y = playerPhysics.Position.y - quant
			rayPos.x = (playerPhysics.Position.y-rayPos.y)*aTan + playerPhysics.Position.x
			offset.y = quant
			offset.x = -offset.y * aTan
		} else if rayAngle == 0 || rayAngle == onePi {
			rayPos.x = playerPhysics.Position.x
			rayPos.y = playerPhysics.Position.y
			dof = maxDof
		}
		for dof < maxDof {
			if rayPos.x > 0 && rayPos.x < float64(mapSize.x) &&
				rayPos.y > 0 && rayPos.y < float64(mapSize.y) {
				red, green, blue, alpha := mapImg.At(int(rayPos.x), int(rayPos.y)).RGBA()
				if (red > 0 || green > 0 || blue > 0) && alpha > 0 {
					ebitenutil.DrawLine(screen, playerPhysics.Position.x*screenScale, playerPhysics.Position.y*screenScale, rayPos.x*screenScale, rayPos.y*screenScale, cRed)
					dof = maxDof
					break
				} else {
					/* next line */
					rayPos.x += playerPhysics.MovePos.x
					rayPos.y += playerPhysics.MovePos.y
					dof += 1
				}
			} else {
				dof = maxDof /* edge of map */
			}
		}
	}

	/* Draw Player */
	ebitenutil.DrawLine(screen,
		playerPhysics.Position.x*screenScale, playerPhysics.Position.y*screenScale,
		playerPhysics.Position.x*screenScale+playerPhysics.MovePos.x*playerLineLen, playerPhysics.Position.y*screenScale+playerPhysics.MovePos.y*playerLineLen,
		cYellow)
	ebitenutil.DrawCircle(screen, playerPhysics.Position.x*screenScale, playerPhysics.Position.y*screenScale, playerCircleCir, cYellow)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
}

func renderMap() {
	if mapDirty {
		mapDirty = false
		mapRender.Fill(color.Transparent)

		/* Draw walls */
		for x := 0; x < mapSize.x; x++ {
			for y := 0; y < mapSize.y; y++ {
				r, g, b, a := mapImg.At(x, y).RGBA()
				if (r > 0 || g > 0 || b > 0) && a > 0 {
					ebitenutil.DrawRect(mapRender, float64(x*screenScale), float64(y*screenScale), screenScale-1, screenScale-1, color.White)
				}
			}
		}
	}
}
