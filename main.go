package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func main() {
	var err error
	playerPhysics.Position.x = 1
	playerPhysics.Position.y = 1

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("GoRaycaster")

	mapImg, _, err = ebitenutil.NewImageFromFile("map1.png")
	if err != nil {
		fmt.Println(err)
		return
	}

	lineImg = ebiten.NewImage(1, 1)
	lineImg.Fill(cYellow)

	wallImg = ebiten.NewImage(1, 1)
	wallImg.Fill(color.White)

	xs, ys := mapImg.Size()
	mapSize.x = float64(xs)
	mapSize.y = float64(ys)

	mapRender = ebiten.NewImage(xs*screenScale, ys*screenScale)

	playerPhysics.MovePos.x = math.Cos(playerPhysics.Rotation)
	playerPhysics.MovePos.y = -math.Sin(playerPhysics.Rotation)

	fmt.Printf("Map size: %v,%v\n", mapSize.x, mapSize.y)

	g := &Game{}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
