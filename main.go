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

	mapRender = ebiten.NewImage(xs*mapScale, ys*mapScale)

	playerPhysics.MovePos.x = math.Cos(playerPhysics.Rotation)
	playerPhysics.MovePos.y = -math.Sin(playerPhysics.Rotation)
	miniMapOffsetX = float64(screenWidth*screenMag) - (float64(xs * mapScale))

	fmt.Printf("Map size: %v,%v\n", mapSize.x, mapSize.y)

	g := &Game{}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func degToRad(deg float64) float64 {
	return fixRad(deg * onePi / 180.0)
}

func fixRad(rad float64) float64 {
	return rad - twoPi*math.Floor((rad+onePi)/twoPi)
}
