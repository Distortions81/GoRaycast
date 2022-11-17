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
	playerPos.x = 1
	playerPos.y = 1

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("GoRaycaster")

	mapImg, _, err = ebitenutil.NewImageFromFile("map1.png")
	if err != nil {
		fmt.Println(err)
		return
	}

	screenCenter.x = screenWidth / 2
	screenCenter.y = screenHeight / 2

	lineImg = ebiten.NewImage(1, 1)
	lineImg.Fill(cYellow)

	xs, ys := mapImg.Size()
	mapSize.x = float64(xs)
	mapSize.y = float64(ys)

	playerPosR.x = math.Cos(playerRot)
	playerPosR.y = -math.Sin(playerRot)

	fmt.Printf("Map size: %v,%v\n", mapSize.x, mapSize.y)

	g := &Game{}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
