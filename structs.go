package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 320
	screenHeight = 240

	mapXSize = 16
	mapYSize = 16

	flatScale = 4
	drawScale = 2

	charMoveSpeed = ( /*blocks per second*/ 8.0 / /*tps*/ 60.0)
)

var cYellow = color.RGBA{0xFF, 0xAA, 0x00, 0xFF}

type Game struct {
	keys []ebiten.Key
}

type xycord struct {
	x float64
	y float64
}
type ixycord struct {
	x int
	y int
}

var flatMap *ebiten.Image
var playerImg *ebiten.Image

var flatSize float64

var playerPos xycord
var oldPlayerPos xycord
