package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth        = 1024
	screenHeight       = 1024
	screenScale        = 64
	playerRotSpeed     = 1.0
	playerForwardSpeed = 10
	twoPi              = math.Pi * 2.0
	onePi              = math.Pi
	halfPi             = math.Pi / 2.0
	quarterPi          = math.Pi / 4.0

	/* long distance run, 2.2 to 2.6m/s */
	/* walking 1.1 to 1.7m/s */
)

var (
	cYellow = color.RGBA{0xFF, 0xAA, 0x00, 0xFF}
	cRed    = color.RGBA{0xFF, 0x00, 0x00, 0xFF}

	mapSize ixycord

	mapImg    *ebiten.Image
	mapRender *ebiten.Image
	mapDirty  bool = true

	playerPhysics pPhysics
)

type pPhysics struct {
	Position xycord
	Rotation float64
	MovePos  xycord
}

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
