package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth   = 320
	screenHeight  = 240
	mapXSize      = 16
	mapYSize      = 16
	flatScale     = 4
	drawScale     = 2
	charMoveSpeed = ( /*blocks per second*/ 4.0 / /*tps*/ 60.0)
)

var GameMap = []uint8{
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
}

var flatMap *ebiten.Image
var playerImg *ebiten.Image
var flatSize float64

type xycord struct {
	x float64
	y float64
}
type ixycord struct {
	x int
	y int
}

var playerPos xycord

type Game struct {
	keys []ebiten.Key
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {

	var op *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	op.GeoM.Scale(flatScale, flatScale)
	op.GeoM.Translate(float64(screenWidth)-flatSize, float64(screenHeight)-flatSize)
	op.Filter = ebiten.FilterNearest
	screen.DrawImage(flatMap, op)

	for _, p := range g.keys {
		if p == ebiten.KeyS {
			if playerPos.y < mapYSize {
				playerPos.y += charMoveSpeed
			}
		}
		if p == ebiten.KeyW {
			if playerPos.y > 0 {
				playerPos.y -= charMoveSpeed
			}
		}
		if p == ebiten.KeyD {
			if playerPos.x < mapXSize {
				playerPos.x += charMoveSpeed
			}
		}
		if p == ebiten.KeyA {
			if playerPos.x > 0 {
				playerPos.x -= charMoveSpeed
			}
		}

		playerImg.Fill(color.Transparent)
		playerImg.Set(int(playerPos.x), int(playerPos.y), color.RGBA{0xff, 0xff, 0x00, 0xff})
	}
	screen.DrawImage(playerImg, op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*drawScale, screenHeight*drawScale)
	ebiten.SetWindowTitle("GoRaycaster")
	//ebiten.SetWindowResizable(true)
	flatMap = ebiten.NewImage(mapXSize, mapYSize)
	playerImg = ebiten.NewImage(mapXSize, mapYSize)
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
