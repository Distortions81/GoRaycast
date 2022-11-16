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
	oldPlayerPos.x = playerPos.x
	oldPlayerPos.y = playerPos.y

	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	for _, p := range g.keys {
		switch p {
		case ebiten.KeyS:
			if playerPos.y < mapYSize-1 {
				playerPos.y += charMoveSpeed
			}
		case ebiten.KeyW:
			if playerPos.y > 0 {
				playerPos.y -= charMoveSpeed
			}

		case ebiten.KeyD:
			if playerPos.x < mapXSize-1 {
				playerPos.x += charMoveSpeed
			}
		case ebiten.KeyA:
			if playerPos.x > 0 {
				playerPos.x -= charMoveSpeed
			}
		}
	}
	if int(oldPlayerPos.x) != int(playerPos.x) || int(oldPlayerPos.y) != int(playerPos.y) {
		playerImg.Fill(color.Transparent)
		playerImg.Set(int(playerPos.x), int(playerPos.y), cYellow)
	}
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {

	var op *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	op.GeoM.Scale(flatScale, flatScale)
	op.GeoM.Translate(float64(screenWidth)-flatSize, float64(screenHeight)-flatSize)
	op.Filter = ebiten.FilterNearest
	screen.DrawImage(flatMap, op)
	screen.DrawImage(playerImg, op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	playerPos.x = 1
	playerPos.y = 1

	ebiten.SetWindowSize(screenWidth*drawScale, screenHeight*drawScale)
	ebiten.SetWindowTitle("GoRaycaster")

	//ebiten.SetWindowResizable(true)
	flatMap = ebiten.NewImage(mapXSize, mapYSize)
	playerImg = ebiten.NewImage(mapXSize, mapYSize)

	/* Init player position */
	playerImg.Fill(color.Transparent)
	playerImg.Set(int(playerPos.x), int(playerPos.y), cYellow)

	updateFlatMap()
	g := &Game{}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func updateFlatMap() {
	flatSize = float64(mapXSize * flatScale)
	for y := 0; y < mapYSize; y++ {
		for x := 0; x < mapXSize; x++ {
			point := y*mapXSize + x
			if GameMap[point] > 0 {
				flatMap.Set(x, y, color.White)
			}
		}
	}
}
