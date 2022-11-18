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
	//rayAngle := playerPhysics.Rotation
	sinRayAngle := math.Sin(playerPhysics.Rotation)
	//cosRayAngle := math.Cos(playerPhysics.Rotation)
	tanRayAngle := math.Tan(playerPhysics.Rotation)
	iTanRayAngle := 1.0 / tanRayAngle

	var rayPos xycord
	var offset xycord

	for rayNum := 0; rayNum < 320; rayNum++ {
		dof := 0

		/* Check Horizontal Lines */
		if sinRayAngle > 0.001 {
			rayPos.y = (math.Floor(playerPhysics.Position.y)/mapScale)*mapScale - 0.0001
			rayPos.x = (playerPhysics.Position.y-rayPos.y)*iTanRayAngle + playerPhysics.Position.x
			offset.y = -mapScale
			offset.x = -offset.y * iTanRayAngle
		} else if sinRayAngle < -0.001 {
			rayPos.y = (math.Floor(playerPhysics.Position.y)/mapScale)*mapScale + mapScale
			rayPos.x = (playerPhysics.Position.y-rayPos.y)*iTanRayAngle + playerPhysics.Position.x
			offset.y = mapScale
			offset.x = -offset.y * iTanRayAngle
		} else {
			rayPos.x = playerPhysics.Position.x
			rayPos.y = playerPhysics.Position.y
			dof = maxDof
		}
		for dof < maxDof {
			if rayPos.x > 0 && rayPos.x < float64(mapSize.x*mapScale) &&
				rayPos.y > 0 && rayPos.y < float64(mapSize.y*mapScale) {
				red, green, blue, alpha := mapImg.At(int(rayPos.x/mapScale), int(rayPos.y/mapScale)).RGBA()
				if (red > 0 || green > 0 || blue > 0) && alpha > 0 {
					ebitenutil.DrawLine(screen, playerPhysics.Position.x, playerPhysics.Position.y, rayPos.x, rayPos.y, cRed)
					dof = maxDof
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
		playerPhysics.Position.x, playerPhysics.Position.y,
		playerPhysics.Position.x+playerPhysics.MovePos.x*playerLineLen, playerPhysics.Position.y+playerPhysics.MovePos.y*playerLineLen,
		cYellow)
	ebitenutil.DrawCircle(screen, playerPhysics.Position.x, playerPhysics.Position.y, playerCircleCir, cYellow)
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
					ebitenutil.DrawRect(mapRender, float64(x*mapScale), float64(y*mapScale), mapScale-1, mapScale-1, color.White)
				} else {
					ebitenutil.DrawRect(mapRender, float64(x*mapScale), float64(y*mapScale), mapScale-1, mapScale-1, cDarkGray)
				}
			}
		}
	}
}
