package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	numRays               = 1
	renderAngle           = 0
	halfRenderAngle       = renderAngle / 2
	screenWidth           = 1280
	screenHeight          = 800
	screenMag             = 1
	mapScale              = 16
	playerLineLen         = 16
	playerCircleCir       = 4
	playerRotSpeed        = 2
	playerForwardSpeedDiv = 1
	threePi               = math.Pi * 3.0
	twoPi                 = math.Pi * 2.0
	onePi                 = math.Pi
	halfPi                = math.Pi / 2.0
	quarterPi             = math.Pi / 4.0

	/* long distance run, 2.2 to 2.6m/s */
	/* walking 1.1 to 1.7m/s */
)

var (
	halfRenderRad  = degToRad(halfRenderAngle)
	rayRads        = degToRad(renderAngle / screenWidth)
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
