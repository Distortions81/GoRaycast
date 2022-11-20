package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func main() {
	var err error

	/* Player starting pos/rot */
	playerPhysics.Position.x = mapScale * 2
	playerPhysics.Position.y = mapScale * 2
	playerPhysics.Rotation = onePi
	angleCalc() /* Update movepos */

	/* Window init */
	ebiten.SetWindowSize(screenWidth*screenMag, screenHeight*screenMag)
	ebiten.SetWindowTitle("GoRaycaster")

	/* Load default test map */
	mapImg, _, err = ebitenutil.NewImageFromFile("map1.png")
	if err != nil {
		fmt.Println(err)
		return
	}

	/* Save size info */
	xs, ys := mapImg.Size()
	mapSize.x = xs
	mapSize.y = ys

	/* Calculate reasonable max depth */
	maxDof = (xs * mapScale) * (ys * mapScale)

	/* Precalc fov values */
	renderFovRad = degToRad(renderFov)
	halfFovRad = fixRad(renderFovRad / 2.0)
	radPerRay = fixRad(renderFovRad / screenWidth)

	/* Minimap setup */
	mapRender = ebiten.NewImage(xs*mapScale, ys*mapScale)
	miniMapOffsetX = float64(screenWidth) - (float64(xs * mapScale))

	fmt.Printf("Map size: %v,%v\n", mapSize.x, mapSize.y)

	/* Start ebiten */
	g := &Game{}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
