package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func main() {
	var err error

	/* Player starting pos/rot */
	playerPhysics.Position.x = mapScale * 2
	playerPhysics.Position.y = mapScale * 2
	playerPhysics.Rotation = 0.0001

	angleCalc() /* Update player  movepos */

	/* Window init */
	ebiten.SetWindowSize(screenWidth*screenMag, screenHeight*screenMag)
	ebiten.SetWindowTitle("GoRaycaster")
	screenSave = ebiten.NewImage(screenWidth, screenHeight)
	screenSave.Fill(color.Black) //Clear screen

	/* Load default test map */
	mapImg, _, err = ebitenutil.NewImageFromFile("map1c.png")
	if err != nil {
		fmt.Println(err)
		return //Exit on error
	}

	/* Load default test title */

	titleImg, _, err = ebitenutil.NewImageFromFile("title.png")
	if err != nil {
		fmt.Println(err)
		return //Exit on error
	}

	/* Meltscreen buffers */
	meltStart = ebiten.NewImage(meltWidth, meltHeight)
	meltBuf = ebiten.NewImage(meltWidth, meltHeight)
	randomizeMelt() //Randomize values
	doMelt = 1      //Start timer

	op := &ebiten.DrawImageOptions{}
	var titleSize ixycord
	titleSize.x, titleSize.y = titleImg.Size()
	op.GeoM.Scale(meltWidth/float64(titleSize.x), meltHeight/float64(titleSize.y))
	meltStart.DrawImage(titleImg, op)
	//meltStart.Fill(cRed)

	/* Save map size info */
	xs, ys := mapImg.Size()
	mapSize.x = xs
	mapSize.y = ys

	/* minimap buffers */
	miniMap = ebiten.NewImage(miniScale*mapSize.x, miniScale*mapSize.y)
	rayImg = ebiten.NewImage(miniScale*mapSize.x, miniScale*mapSize.y)

	/* Calculate reasonable max depth */
	maxDof = ((xs + ys) / 2) * 2

	/* Precalc fov values */
	renderFovRad = degToRad(renderFov)
	halfFovRad = fixRad(renderFovRad / 2.0)
	radPerRay = fixRad(renderFovRad / screenWidth)

	fmt.Printf("Map size: %v,%v\n", mapSize.x, mapSize.y)

	/* Start ebiten */
	g := &Game{}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
