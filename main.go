package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {
	oldPlayerPos = playerPos

	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	pad := float64(1.0 / mapScale)
	for _, p := range g.keys {
		switch p {
		case ebiten.KeyS:
			if playerPos.y < float64(mapSize.y)-pad {
				playerPos.y += charMoveSpeed
			}

		case ebiten.KeyW:
			if playerPos.y > 0 {
				playerPos.y -= charMoveSpeed
			}

		case ebiten.KeyD:
			if playerPos.x < float64(mapSize.x)-pad {
				playerPos.x += charMoveSpeed
			}
		case ebiten.KeyA:
			if playerPos.x > 0 {
				playerPos.x -= charMoveSpeed
			}
		}
	}
	if oldPlayerPos != playerPos {
		playerImg.Fill(color.Transparent)
		playerImg.Set(int(playerPos.x*mapScale), int(playerPos.y*mapScale), cYellow)
	}
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {

	var op *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	op.GeoM.Scale(mapScale, mapScale)
	op.GeoM.Translate(float64(screenWidth-mapSize.x*mapScale), float64(screenHeight-mapSize.y*mapScale))
	op.Filter = ebiten.FilterNearest
	screen.DrawImage(mapImg, op)

	op.GeoM.Reset()
	op.GeoM.Translate(float64(screenWidth-mapSize.x*mapScale), float64(screenHeight-mapSize.y*mapScale))
	op.Filter = ebiten.FilterNearest
	screen.DrawImage(playerImg, op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	var err error
	playerPos.x = 1
	playerPos.y = 1

	ebiten.SetWindowSize(screenWidth*screenScale, screenHeight*screenScale)
	ebiten.SetWindowTitle("GoRaycaster")

	mapImg, _, err = ebitenutil.NewImageFromFile("map1.png")
	if err != nil {
		fmt.Println(err)
		return
	}

	mapSize.x, mapSize.y = mapImg.Size()
	playerImg = ebiten.NewImage(mapSize.x*mapScale, mapSize.y*mapScale)
	fmt.Printf("Map size: %v,%v\n", mapSize.x, mapSize.y)

	g := &Game{}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
