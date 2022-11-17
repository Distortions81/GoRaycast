package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth     = 1280 / 2
	screenHeight    = 720 / 2
	screenScale     = 5
	playerMoveSpeed = 5.0
	twoPi           = math.Pi * 2.0
	onePi           = math.Pi
	halfPi          = math.Pi / 2.0
	quarterPi       = math.Pi / 4.0

	/* long distance run, 2.2 to 2.6m/s */
	/* walking 1.1 to 1.7m/s */
)

var (
	playerRotSpeed = 0.166
	screenCenter   xycord

	cYellow = color.RGBA{0xFF, 0xAA, 0x00, 0xFF}
	cRed    = color.RGBA{0xFF, 0x00, 0x00, 0xFF}

	mapSize xycord

	mapImg  *ebiten.Image
	lineImg *ebiten.Image

	playerPos  xycord
	playerRot  float64
	playerPosR xycord
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
