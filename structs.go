package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	maxDist               = 1000000.0
	renderFov             = 60
	screenWidth           = 1920
	screenHeight          = 1080
	screenMag             = 1
	mapScale              = 32
	playerLineLen         = 16
	playerCircleCir       = 4
	playerRotSpeed        = 2
	playerForwardSpeedDiv = 2
	threePi               = math.Pi * 3.0
	twoPi                 = math.Pi * 2.0
	onePi                 = math.Pi
	halfPi                = math.Pi / 2.0
	quarterPi             = math.Pi / 4.0

	/* long distance run, 2.2 to 2.6m/s */
	/* walking 1.1 to 1.7m/s */
)

var (
	renderFovRad float64
	halfFovRad   float64
	radPerRay    float64

	miniMapOffsetX float64
	cDarkGray      = color.RGBA{0x20, 0x20, 0x20, 0xFF}
	cYellow        = color.RGBA{0xFF, 0xAA, 0x00, 0xFF}
	cRed           = color.RGBA{0xFF, 0x00, 0x00, 0xFF}
	cGreen         = color.RGBA{0x00, 0xFF, 0x00, 0xFF}

	mapSize ixycord

	mapImg    *ebiten.Image
	mapRender *ebiten.Image
	mapDirty  bool = true
	maxDof    int

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
