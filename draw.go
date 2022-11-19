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
	op.GeoM.Translate(miniMapOffsetX, 0)
	screen.DrawImage(mapRender, op)

	/* Draw rays */

	var vrayPos xycord
	var hrayPos xycord
	var offset xycord
	var movePos xycord
	rayAngle := fixRad(playerPhysics.Rotation + halfFovRad)

	for rayNum := 0; rayNum < numRays; rayNum++ {
		movePos.x = math.Cos(rayAngle)  // opposite
		movePos.y = -math.Sin(rayAngle) // adjacent

		dof := 0
		sinRayAngle := math.Sin(rayAngle)
		tanRayAngle := math.Tan(rayAngle)
		iTanRayAngle := 1.0 / tanRayAngle

		hDist := maxDist
		vDist := maxDist

		/* Check Vertical Lines */
		if sinRayAngle > 0.001 { //Look left
			vrayPos.x = (math.Floor(playerPhysics.Position.x)/mapScale)*mapScale + mapScale
			vrayPos.y = (playerPhysics.Position.x-vrayPos.x)*tanRayAngle + playerPhysics.Position.y
			offset.x = mapScale
			offset.y = -offset.x * tanRayAngle
		} else if sinRayAngle < -0.001 { //Look right
			vrayPos.x = (math.Floor(playerPhysics.Position.x)/mapScale)*mapScale - 0.0001
			vrayPos.y = (playerPhysics.Position.x-vrayPos.x)*tanRayAngle + playerPhysics.Position.y
			offset.x = -mapScale
			offset.y = -offset.x * tanRayAngle
		} else {
			vrayPos.x = playerPhysics.Position.x
			vrayPos.y = playerPhysics.Position.y
			dof = maxDof
		}
		for dof < maxDof {
			if vrayPos.x >= 0 && vrayPos.x <= float64(mapSize.x*mapScale) &&
				vrayPos.y >= 0 && vrayPos.y <= float64(mapSize.y*mapScale) {
				red, green, blue, alpha := mapImg.At(int(vrayPos.x/mapScale), int(vrayPos.y/mapScale)).RGBA()
				if (red > 0 || green > 0 || blue > 0) && alpha > 0 {
					vDist = distance(playerPhysics.Position, vrayPos)
					dof = maxDof
				} else {
					/* next line */
					vrayPos.x += movePos.x
					vrayPos.y += movePos.y

					dof += 1
				}
			} else {
				dof = maxDof /* edge of map */
			}
		}

		/* Check Horizontal Lines */
		dof = 0
		if sinRayAngle > 0.001 { //Look north
			hrayPos.y = (math.Floor(playerPhysics.Position.y)/mapScale)*mapScale - 0.0001
			hrayPos.x = (playerPhysics.Position.y-hrayPos.y)*iTanRayAngle + playerPhysics.Position.x
			offset.y = -mapScale
			offset.x = -offset.y * iTanRayAngle
		} else if sinRayAngle < -0.001 { //Look south
			hrayPos.y = (math.Floor(playerPhysics.Position.y)/mapScale)*mapScale + mapScale
			hrayPos.x = (playerPhysics.Position.y-hrayPos.y)*iTanRayAngle + playerPhysics.Position.x
			offset.y = mapScale
			offset.x = -offset.y * iTanRayAngle
		} else {
			hrayPos.x = playerPhysics.Position.x
			hrayPos.y = playerPhysics.Position.y
			dof = maxDof
		}
		for dof < maxDof {
			if hrayPos.x >= 0 && hrayPos.x <= float64(mapSize.x*mapScale) &&
				hrayPos.y >= 0 && hrayPos.y <= float64(mapSize.y*mapScale) {
				red, green, blue, alpha := mapImg.At(int(hrayPos.x/mapScale), int(hrayPos.y/mapScale)).RGBA()
				if (red > 0 || green > 0 || blue > 0) && alpha > 0 {
					hDist = distance(playerPhysics.Position, hrayPos)
					dof = maxDof
				} else {
					/* next line */
					hrayPos.x += movePos.x
					hrayPos.y += movePos.y

					dof += 1
				}
			} else {
				dof = maxDof /* edge of map */
			}
		}

		//Use shortest vector, if any found
		if hDist < maxDist || vDist < maxDist {
			if hDist < vDist {
				ebitenutil.DrawLine(screen, miniMapOffsetX+playerPhysics.Position.x, playerPhysics.Position.y, miniMapOffsetX+hrayPos.x, hrayPos.y, cRed)
			} else {
				ebitenutil.DrawLine(screen, miniMapOffsetX+playerPhysics.Position.x, playerPhysics.Position.y, miniMapOffsetX+vrayPos.x, vrayPos.y, cRed)
			}
		}

		rayAngle = fixRad(rayAngle - radPerRay)
	}

	/* Draw Player */
	ebitenutil.DrawLine(screen,
		miniMapOffsetX+playerPhysics.Position.x, playerPhysics.Position.y,
		miniMapOffsetX+playerPhysics.Position.x+playerPhysics.MovePos.x*playerLineLen, playerPhysics.Position.y+playerPhysics.MovePos.y*playerLineLen,
		cYellow)
	ebitenutil.DrawCircle(screen, miniMapOffsetX+playerPhysics.Position.x, playerPhysics.Position.y, playerCircleCir, cYellow)
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
					ebitenutil.DrawRect(mapRender, float64(x*mapScale), float64(y*mapScale), mapScale-1, mapScale-1, color.Black)
				} else {
					ebitenutil.DrawRect(mapRender, float64(x*mapScale), float64(y*mapScale), mapScale-1, mapScale-1, cDarkGray)
				}
			}
		}
	}
}
