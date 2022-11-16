package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1280 / 2
	screenHeight = 720 / 2
	screenScale  = 2
	mapScale     = 4

	/* long distance run, 2.2 to 2.6m/s */
	/* walking 1.1 to 1.7m/s */
	charMoveSpeed = ( /*blocks per second*/ 8.0 / /*tps*/ 60.0)
)

var (
	cYellow = color.RGBA{0xFF, 0xAA, 0x00, 0xFF}
	cRed    = color.RGBA{0xFF, 0x00, 0x00, 0xFF}

	mapSize ixycord
)

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

var mapImg *ebiten.Image
var playerImg *ebiten.Image

var playerPos xycord
var oldPlayerPos xycord
