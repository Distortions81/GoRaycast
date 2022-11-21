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
	/* Process input */
	g.Update()

	var verticalRayPosition xycord
	var horizontalRayPosition xycord
	var offset xycord

	/* Move ray counter-clockwise from center, half our FOV */
	rayAngle := fixRad(playerPhysics.Rotation + halfFovRad)

	/* Cast rays */
	for rayNum := 0; rayNum < screenWidth; rayNum++ {

		/* Reset depth */
		currentDepth := 0

		/* Reset offsets */
		offset.x = 0
		offset.y = 0

		/* Reset ray distance */
		horizontalDistance := maxDist
		verticalDistance := maxDist
		finalDistance := 0.0

		/* Precalc commonly used conversions */
		cosRayAngle := math.Cos(rayAngle)
		sinRayAngle := math.Sin(rayAngle)
		tanRayAngle := math.Tan(rayAngle)
		iTanRayAngle := 1.0 / tanRayAngle

		/* Check Vertical Lines */
		if cosRayAngle > 0.001 { //Look left
			verticalRayPosition.x = math.Floor(playerPhysics.Position.x/mapScale)*mapScale + mapScale
			verticalRayPosition.y = (playerPhysics.Position.x-verticalRayPosition.x)*tanRayAngle + playerPhysics.Position.y
			offset.x = mapScale
			offset.y = -offset.x * tanRayAngle
		} else if cosRayAngle < -0.001 { //Look right
			verticalRayPosition.x = math.Floor(playerPhysics.Position.x/mapScale)*mapScale - 0.0001
			verticalRayPosition.y = (playerPhysics.Position.x-verticalRayPosition.x)*tanRayAngle + playerPhysics.Position.y
			offset.x = -mapScale
			offset.y = -offset.x * tanRayAngle
		} else {
			verticalRayPosition.x = playerPhysics.Position.x
			verticalRayPosition.y = playerPhysics.Position.y
			currentDepth = maxDof // Skip loop
		}
		for currentDepth < maxDof {
			if verticalRayPosition.x >= 0 && verticalRayPosition.x <= float64(mapSize.x*mapScale) &&
				verticalRayPosition.y >= 0 && verticalRayPosition.y <= float64(mapSize.y*mapScale) {

				/* Check if there is a wall here */
				red, _, _, _ := mapImg.At(int(verticalRayPosition.x/mapScale), int(verticalRayPosition.y/mapScale)).RGBA()
				if red > 0 {
					/* Calc distance, save, exit */
					verticalDistance = distance(playerPhysics.Position, verticalRayPosition)
					currentDepth = maxDof
				} else {
					/* advance down the ray angle */
					verticalRayPosition.x += offset.x
					verticalRayPosition.y += offset.y

					currentDepth += 1
				}
			} else {
				currentDepth = maxDof /* edge of map, exit loop*/
			}
		}

		/* Reset depth */
		currentDepth = 0

		/* Check Horizontal Lines */
		if sinRayAngle > 0.001 { //Look north
			horizontalRayPosition.y = math.Floor(playerPhysics.Position.y/mapScale)*mapScale - 0.0001
			horizontalRayPosition.x = (playerPhysics.Position.y-horizontalRayPosition.y)*iTanRayAngle + playerPhysics.Position.x
			offset.y = -mapScale
			offset.x = -offset.y * iTanRayAngle
		} else if sinRayAngle < -0.001 { //Look south
			horizontalRayPosition.y = math.Floor(playerPhysics.Position.y/mapScale)*mapScale + mapScale
			horizontalRayPosition.x = (playerPhysics.Position.y-horizontalRayPosition.y)*iTanRayAngle + playerPhysics.Position.x
			offset.y = mapScale
			offset.x = -offset.y * iTanRayAngle
		} else {
			horizontalRayPosition.x = playerPhysics.Position.x
			horizontalRayPosition.y = playerPhysics.Position.y
			currentDepth = maxDof
		}
		for currentDepth < maxDof {
			if horizontalRayPosition.x >= 0 && horizontalRayPosition.x <= float64(mapSize.x*mapScale) &&
				horizontalRayPosition.y >= 0 && horizontalRayPosition.y <= float64(mapSize.y*mapScale) {

				/* Check if there is a wall here */
				red, _, _, _ := mapImg.At(int(horizontalRayPosition.x/mapScale), int(horizontalRayPosition.y/mapScale)).RGBA()
				if red > 0 {
					/* Calc distance, save, exit */
					horizontalDistance = distance(playerPhysics.Position, horizontalRayPosition)
					currentDepth = maxDof
				} else {
					/* dvance down the ay angle */
					horizontalRayPosition.x += offset.x
					horizontalRayPosition.y += offset.y

					currentDepth += 1
				}
			} else {
				currentDepth = maxDof /* edge of map, exit loop */
			}
		}

		//Use shortest vector, if any found
		if horizontalDistance < maxDist || verticalDistance < maxDist {
			if horizontalDistance < verticalDistance {
				finalDistance = horizontalDistance
			} else {
				finalDistance = verticalDistance
			}
		}

		/* Draw rays */
		if finalDistance < maxDist {
			lh := (float64(mapSize.y) * screenHeight) / finalDistance
			bright := uint8(float64(mapSize.y*255) / finalDistance)
			ebitenutil.DrawRect(screen, float64(rayNum), (screenHeight/2.0)-(lh/2.0), 1, lh, color.RGBA{bright, bright, bright, 0xFF})
		}

		/* Advance ray angle */
		rayAngle = fixRad(rayAngle - radPerRay)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
}
