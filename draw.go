package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	frameNumber++

	/* Process input */
	g.processInput(screen)

	var s *ebiten.Image
	if doMelt < 0 || meltQuit {
		screenSave.Fill(color.Black)
		s = screenSave
	} else {
		s = screen
	}

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

		var verticalColor color.Color
		var horizontalColor color.Color
		var finalColor color.Color
		finalColor = color.White

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
				verticalColor = mapImg.At(int(verticalRayPosition.x/mapScale), int(verticalRayPosition.y/mapScale))
				if r, g, b, _ := verticalColor.RGBA(); r > 0 || g > 0 || b > 0 {
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
				horizontalColor = mapImg.At(int(horizontalRayPosition.x/mapScale), int(horizontalRayPosition.y/mapScale))
				if r, g, b, _ := horizontalColor.RGBA(); r > 0 || g > 0 || b > 0 {
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
				finalColor = horizontalColor
			} else {
				finalDistance = verticalDistance
				finalColor = verticalColor
			}
		}

		/* Draw rays */
		if finalDistance < maxDist {
			lh := (float64(mapSize.y) * screenHeight) / finalDistance
			r, g, b, _ := finalColor.RGBA()
			d := (float64(mapSize.y+mapSize.x) / (finalDistance))
			if d < 0 {
				d = 0
			} else if d > 1 {
				d = 1
			}
			red := uint8(((float64(r) / 255.0) * d))
			green := uint8(((float64(g) / 255.0) * d))
			blue := uint8(((float64(b) / 255.0) * d))
			ebitenutil.DrawRect(s, float64(rayNum), (screenHeight/2.0)-(lh/2.0), 1, lh, color.RGBA{red, green, blue, 0xFF})
		}

		/* Advance ray angle */
		rayAngle = fixRad(rayAngle - radPerRay)
	}
	ebitenutil.DebugPrint(s, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))

	if doMelt < 0 {
		doMelt = meltFrames
		op := &ebiten.DrawImageOptions{}
		var scale xycord
		scale.x = screenWidth / meltWidth
		scale.y = screenHeight / meltHeight
		op.GeoM.Scale(1.0/scale.x, 1.0/scale.y)
		meltStart.DrawImage(screenSave, op)
		screen.DrawImage(screenSave, nil)

		randomizeMelt()
	}

	if doMelt > 0 {
		doMelt--
		op := &ebiten.DrawImageOptions{}

		meltBuf.Fill(color.Transparent)
		for i := 0; i < meltWidth; i++ {
			op.GeoM.Reset()
			offset := meltOffsets[i]
			fNum := (meltFrames - doMelt) * meltSpeed
			newOff := 0
			if fNum > offset {
				newOff = fNum - offset
			}
			op.GeoM.Translate(float64(i), float64(newOff))
			meltBuf.DrawImage(meltStart.SubImage(image.Rect(i, 0, i+1, meltHeight)).(*ebiten.Image), op)
		}

		/*Draw to screen */
		var meltScale xycord
		meltScale.x = screenWidth / meltWidth
		meltScale.y = screenHeight / meltHeight
		op.GeoM.Reset()
		op.GeoM.Scale(meltScale.x, meltScale.y)
		screen.DrawImage(meltBuf, op)

		if doMelt == 0 && meltQuit {
			os.Exit(1)
		}
	}
}

func randomizeMelt() {
	for i := 0; i < meltWidth; i++ {
		meltOffsets[i] = rand.Intn(meltAmount)
	}
}
