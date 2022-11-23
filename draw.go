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

	/* Process user input */
	g.processInput(screen)

	var s *ebiten.Image //Pointer, so we can screen cap if we want

	/* If we are set to melt, or to quit game, swap output pointer */
	if doMelt < 0 || meltQuit {
		//screenSave.Fill(cRed)
		s = screenSave
	} else {
		s = screen
	}

	/* Move ray counter-clockwise from center, half our FOV */
	rayAngle := fixRad(playerPhysics.Rotation + halfFovRad)

	/* Cast rays */
	for rayNum := 0; rayNum < screenWidth; rayNum++ {

		/* Reset depth */
		currentDepth := 0

		/* Switches for horizontal or vertical lines */
		var verticalRayPosition xycord
		var horizontalRayPosition xycord
		var offset xycord
		var verticalColor color.Color
		var horizontalColor color.Color
		var finalColor color.Color
		var wallWasHorizontal = false //Used for shading
		finalColor = color.White

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
		if cosRayAngle > 0.001 { /* Looking left */
			verticalRayPosition.x = math.Floor(playerPhysics.Position.x/mapScale)*mapScale + mapScale
			verticalRayPosition.y = (playerPhysics.Position.x-verticalRayPosition.x)*tanRayAngle + playerPhysics.Position.y
			offset.x = mapScale
			offset.y = -offset.x * tanRayAngle
		} else if cosRayAngle < -0.001 { /* Looking right */
			verticalRayPosition.x = math.Floor(playerPhysics.Position.x/mapScale)*mapScale - 0.0001
			verticalRayPosition.y = (playerPhysics.Position.x-verticalRayPosition.x)*tanRayAngle + playerPhysics.Position.y
			offset.x = -mapScale
			offset.y = -offset.x * tanRayAngle
		} else {
			verticalRayPosition.x = playerPhysics.Position.x
			verticalRayPosition.y = playerPhysics.Position.y
			currentDepth = maxDof // Skip loop
		}

		/* Look for a wall */
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
				currentDepth = maxDof /* past edge of map, exit loop */
			}
		}

		/* Reset depth */
		currentDepth = 0

		/* Check Horizontal Lines */
		if sinRayAngle > 0.001 { /* Looking north */
			horizontalRayPosition.y = math.Floor(playerPhysics.Position.y/mapScale)*mapScale - 0.0001
			horizontalRayPosition.x = (playerPhysics.Position.y-horizontalRayPosition.y)*iTanRayAngle + playerPhysics.Position.x
			offset.y = -mapScale
			offset.x = -offset.y * iTanRayAngle
		} else if sinRayAngle < -0.001 { /* Looking south */
			horizontalRayPosition.y = math.Floor(playerPhysics.Position.y/mapScale)*mapScale + mapScale
			horizontalRayPosition.x = (playerPhysics.Position.y-horizontalRayPosition.y)*iTanRayAngle + playerPhysics.Position.x
			offset.y = mapScale
			offset.x = -offset.y * iTanRayAngle
		} else {
			horizontalRayPosition.x = playerPhysics.Position.x
			horizontalRayPosition.y = playerPhysics.Position.y
			currentDepth = maxDof // Skip loop
		}

		/* Look for a wall */
		for currentDepth < maxDof {
			/* Check if position is on the map */
			if horizontalRayPosition.x >= 0 && horizontalRayPosition.x <= float64(mapSize.x*mapScale) &&
				horizontalRayPosition.y >= 0 && horizontalRayPosition.y <= float64(mapSize.y*mapScale) {

				/* Check if there is a wall here, black is empty */
				horizontalColor = mapImg.At(int(horizontalRayPosition.x/mapScale), int(horizontalRayPosition.y/mapScale))

				if r, g, b, _ := horizontalColor.RGBA(); r > 0 || g > 0 || b > 0 {
					/* Calc distance, save, exit */
					horizontalDistance = distance(playerPhysics.Position, horizontalRayPosition)
					currentDepth = maxDof
				} else {
					/* advance down the ray angle */
					horizontalRayPosition.x += offset.x
					horizontalRayPosition.y += offset.y

					currentDepth += 1
				}
			} else {
				currentDepth = maxDof /* past edge of map, exit loop */
			}
		}

		//Use shortest vector, if any found
		if horizontalDistance < maxDist || verticalDistance < maxDist {
			if horizontalDistance < verticalDistance {
				finalDistance = horizontalDistance
				finalColor = horizontalColor
				wallWasHorizontal = true //For shading
			} else {
				finalDistance = verticalDistance
				finalColor = verticalColor
			}
		}

		/* Draw rays */
		if finalDistance < maxDist {

			//Ray length, scaled to map size
			lh := (float64(mapSize.y) * screenHeight) / finalDistance

			/* Color of map block */
			r, g, b, _ := finalColor.RGBA()
			d := (float64(mapSize.y+mapSize.x) / (finalDistance))

			/* Clip brightness */
			if d < maxShadow {
				d = maxShadow
			} else if d > 1 {
				d = 1
			}

			/* Made horizontal walls darker, faux shading */
			shade := 256.0
			if wallWasHorizontal {
				shade = 256.0 * (dirShading)
			}

			red := uint8(((float64(r) / shade) * d))
			green := uint8(((float64(g) / shade) * d))
			blue := uint8(((float64(b) / shade) * d))

			/* Draw the vertical line! */
			ebitenutil.DrawRect(s, float64(rayNum), (screenHeight/2.0)-(lh/2.0), 1, lh, color.RGBA{red, green, blue, 0xFF})
		}

		/* Advance ray angle */
		rayAngle = fixRad(rayAngle - radPerRay)
	}
	ebitenutil.DebugPrint(s, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))

	/* Melt started, grab screen */
	if doMelt < 0 {
		doMelt = meltFrames //start timer

		op := &ebiten.DrawImageOptions{}
		var scale xycord

		/* Scale image down */
		scale.x = screenWidth / meltWidth
		scale.y = screenHeight / meltHeight
		op.GeoM.Scale(1.0/scale.x, 1.0/scale.y)

		/* Save to buffer, and draw to screen */
		meltStart.DrawImage(screenSave, op)
		screen.DrawImage(screenSave, nil)

		/* Randomize values */
		randomizeMelt()
	}

	/* Draw melt */
	if doMelt > 0 {
		doMelt--
		op := &ebiten.DrawImageOptions{}

		/* Clear buffer */
		meltBuf.Fill(color.Transparent)

		/* Loop through each column */
		for i := 0; i < meltWidth; i++ {
			op.GeoM.Reset()
			offset := meltOffsets[i]

			/* How much to move each column */
			fNum := (meltFrames - doMelt) * meltSpeed

			/* Don't start moving until we pass our offset */
			newOff := 0
			if fNum > offset {
				newOff = fNum - offset
			}

			/* Move the line down, and only draw one line at a time. */
			op.GeoM.Translate(float64(i), float64(newOff))

			/* Draw to buffer */
			meltBuf.DrawImage(meltStart.SubImage(image.Rect(i, 0, i+1, meltHeight)).(*ebiten.Image), op)
		}

		/* Draw to screen */
		var meltScale xycord
		meltScale.x = screenWidth / meltWidth
		meltScale.y = screenHeight / meltHeight
		op.GeoM.Reset()
		op.GeoM.Scale(meltScale.x, meltScale.y)
		screen.DrawImage(meltBuf, op)

		/* Marked to exit game, quit */
		if doMelt == 0 && meltQuit {
			os.Exit(1)
		}
	}
}

/* Random offsets for melting */
func randomizeMelt() {
	for i := 0; i < meltWidth; i++ {
		meltOffsets[i] = rand.Intn(meltAmount)
	}

	//Smoothing
	var a, b, c, d, e, f, g int
	for i := 0; i < meltWidth; i++ {
		if i > 2 {
			a = meltOffsets[i-3]
		}
		if i > 1 {
			b = meltOffsets[i-2]
		}
		if i > 0 {
			c = meltOffsets[i-1]
		}
		if i >= 0 {
			d = meltOffsets[i]
		}
		if i < meltWidth-1 {
			e = meltOffsets[i+1]
		}
		if i < meltWidth-2 {
			f = meltOffsets[i+2]
		}
		if i < meltWidth-3 {
			g = meltOffsets[i+3]
		}
		meltOffsets[i] = (a + b + c + d + e + f + g) / 7.0
	}
}
