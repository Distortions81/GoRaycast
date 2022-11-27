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
	if doMelt == 0 {
		g.processInput(screen)
	}

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

	rayImg.Fill(color.Transparent)

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
		var finalRayPosition xycord

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
			offset.x = 0
			offset.y = 0
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
			offset.x = 0
			offset.y = 0
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

		imgRow := 0
		//Use shortest vector, if any found
		if horizontalDistance < maxDist || verticalDistance < maxDist {
			if horizontalDistance < verticalDistance {
				finalDistance = horizontalDistance
				finalColor = horizontalColor
				wallWasHorizontal = true //For shading
				finalRayPosition = horizontalRayPosition
				imgRow = int((float64(finalRayPosition.x/mapScale - math.Floor(finalRayPosition.x/mapScale))) * float64(wallSize.x))
			} else {
				finalDistance = verticalDistance
				finalColor = verticalColor
				finalRayPosition = verticalRayPosition
				imgRow = int((float64(finalRayPosition.y/mapScale - math.Floor(finalRayPosition.y/mapScale))) * float64(wallSize.x))
			}
		}

		/* Draw rays */
		if finalDistance < maxDist {

			/* Color of map block */
			r, g, b, _ := finalColor.RGBA()
			d := shadowBase / (math.Sqrt(math.Pow(finalDistance+distanceOffset, shadowExp)*FalloffRatio) / shadowDistance)

			/* Made horizontal walls darker, faux shading */
			shade := normalShading
			if wallWasHorizontal {
				shade = dirShading
			}
			if d > shadowClip {
				d = shadowClip
			}
			d = d * shade

			rc := (float64(r) / 0xFFFF)
			gc := (float64(g) / 0xFFFF)
			bc := (float64(b) / 0xFFFF)

			//Draw ray lines here, to rayImg
			if rayNum%miniRayMod == 0 {
				ebitenutil.DrawLine(rayImg, (finalRayPosition.x/mapScale)*miniScale, (finalRayPosition.y/mapScale)*miniScale, (playerPhysics.Position.x/mapScale)*miniScale, (playerPhysics.Position.y/mapScale)*miniScale, cRay)
			}

			/* Draw textures */
			op := &ebiten.DrawImageOptions{}
			//op.Filter = ebiten.FilterLinear
			op.GeoM.Translate(0, -float64(wallSize.y/2.0))
			op.GeoM.Scale(1, (mapScale/wallHeightRatio)/finalDistance)
			op.GeoM.Translate(float64(rayNum), (screenHeight / 2.0))
			op.ColorM.Scale(rc, gc, bc, 1.0) //Apply wall color
			op.ColorM.Scale(d, d, d, 1.0)    //Apply shading and depth
			s.DrawImage(wallImg.SubImage(image.Rect(imgRow, 0, imgRow+1, wallSize.y)).(*ebiten.Image), op)
		}

		/* Advance ray angle */
		rayAngle = fixRad(rayAngle - radPerRay)
	}

	/* MiniMap */
	miniMap.Fill(color.Black)
	/* Draw walls -- TODO cache me */
	for x := 0; x < int(mapSize.x); x++ {
		for y := 0; y < int(mapSize.y); y++ {
			r, g, b, _ := mapImg.At(x, y).RGBA()
			ebitenutil.DrawRect(miniMap, float64(x*miniScale), float64(y*miniScale), miniScale-1, miniScale-1, color.RGBA{uint8(r), uint8(g), uint8(b), 0xFF})
		}
	}

	/* Draw rays */
	op := &ebiten.DrawImageOptions{}
	miniMap.DrawImage(rayImg, op)
	/* Draw player */
	ebitenutil.DrawCircle(miniMap, ((playerPhysics.Position.x / mapScale) * miniScale), ((playerPhysics.Position.y / mapScale) * miniScale), 4, cYellow)
	/* Draw to screen */
	op.GeoM.Translate(screenWidth-float64((mapSize.x+1)*miniScale), miniScale)
	op.ColorM.Scale(1, 1, 1, 0.5)
	screen.DrawImage(miniMap, op)

	ebitenutil.DebugPrint(s, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))

	if doMelt < 0 {
		doMelt = 1 //start timer

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
	} else if doMelt > 0 {
		op := &ebiten.DrawImageOptions{}
		if meltDelay > 0 {
			meltDelay--
		} else {
			doMelt++
		}

		/* Clear buffer */
		meltBuf.Fill(color.Transparent)
		isDone := true

		/* Loop through each column */
		for i := 0; i < meltWidth; i++ {
			d := i
			op.GeoM.Reset()
			offset := meltOffsets[i]

			/* Don't start moving until we pass our offset */
			newOff := 0
			if doMelt-2 > offset {
				newOff = (doMelt - 2) - offset
			}
			newOff *= meltSpeed

			/* Move the line down, and only draw one line at a time. */
			op.GeoM.Translate(float64(d), float64(newOff))

			/* Draw to buffer */
			meltBuf.DrawImage(meltStart.SubImage(image.Rect(d, 0, d+1, meltHeight)).(*ebiten.Image), op)
			if newOff < meltHeight+10 {
				isDone = false
			}
		}
		if isDone {
			//fmt.Printf("melt done: %v\n", doMelt)
			doMelt = 0
			meltDelay = 0
		}

		/* Draw to screen */
		var meltScale xycord
		meltScale.x = screenWidth / meltWidth
		meltScale.y = screenHeight / meltHeight
		op.GeoM.Reset()
		op.GeoM.Scale(meltScale.x, meltScale.y)
		screen.DrawImage(meltBuf, op)
		//time.Sleep(time.Millisecond * 500)

		/* Marked to exit game, quit */
		if doMelt == 0 && meltQuit {
			os.Exit(1)
		}
	}
}

/* Random offsets for melting */
func randomizeMelt() {
	r := 0
	meltOffsets[0] = rand.Intn(255) % 16
	for i := 1; i < meltWidth; i++ {
		r = (rand.Intn(255) % 3) - 1
		meltOffsets[i] = meltOffsets[i-1] + r

		if meltOffsets[i] < 0 {
			meltOffsets[i] = 1
		} else if meltOffsets[i] > 15 {
			meltOffsets[i] = 15
		}
		//fmt.Printf("%v, ", meltOffsets[i])
	}
}
