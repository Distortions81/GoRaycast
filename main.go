package main

import (
	"fmt"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func main() {
	var err error
	playerPhysics.Position.x = mapScale * 2
	playerPhysics.Position.y = mapScale * 2
	playerPhysics.Rotation = onePi
	angleCalc()

	ebiten.SetWindowSize(screenWidth*screenMag, screenHeight*screenMag)
	ebiten.SetWindowTitle("GoRaycaster")

	mapImg, _, err = ebitenutil.NewImageFromFile("map1.png")
	if err != nil {
		fmt.Println(err)
		return
	}

	xs, ys := mapImg.Size()
	mapSize.x = xs
	mapSize.y = ys

	maxDof = (xs * mapScale) * (ys * mapScale)

	renderFovRad = degToRad(renderFov)
	halfFovRad = fixRad(renderFovRad / 2.0)
	radPerRay = fixRad(renderFovRad / screenWidth)

	mapRender = ebiten.NewImage(xs*mapScale, ys*mapScale)

	playerPhysics.MovePos.x = math.Cos(playerPhysics.Rotation)
	playerPhysics.MovePos.y = -math.Sin(playerPhysics.Rotation)
	miniMapOffsetX = float64(screenWidth) - (float64(xs * mapScale))

	fmt.Printf("Map size: %v,%v\n", mapSize.x, mapSize.y)

	g := &Game{}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func fixRad(rad float64) float64 {
	for rad < 0 {
		rad += twoPi
	}
	for rad >= twoPi {
		rad -= twoPi
	}
	return rad
}

func distance(a, b xycord) float64 {
	return math.Hypot(float64(a.x-b.x), float64(a.y-b.y))
}

func degToRad(deg float64) float64 {
	return fixRad(deg * onePi / 180.0)
}
