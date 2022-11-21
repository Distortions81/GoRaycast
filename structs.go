package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	maxDist      = 1000000.0
	renderFov    = 90
	screenWidth  = 3840
	screenHeight = 2160
	screenMag    = 1
	mapScale     = 32
	maxShadow    = 0.01
	dirShading   = 2.0 //2.0 would be 50% darker on horizontal walls

	playerLineLen         = 16
	playerCircleCir       = 4
	playerRotSpeed        = 1
	playerForwardSpeedDiv = 0.5

	meltWidth  = screenWidth / 6
	meltHeight = screenHeight / 6
	meltFrames = meltHeight + meltAmount
	meltSpeed  = 3
	meltAmount = 24

	threePi   = math.Pi * 3.0
	twoPi     = math.Pi * 2.0
	onePi     = math.Pi
	halfPi    = math.Pi / 2.0
	quarterPi = math.Pi / 4.0

	/* long distance run, 2.2 to 2.6m/s */
	/* walking 1.1 to 1.7m/s */
)

var (
	frameNumber  uint64
	doMelt       int
	renderFovRad float64
	halfFovRad   float64
	radPerRay    float64

	cDarkGray = color.RGBA{0x20, 0x20, 0x20, 0xFF}
	cYellow   = color.RGBA{0xFF, 0xAA, 0x00, 0xFF}
	cRed      = color.RGBA{0xFF, 0x00, 0x00, 0xFF}
	cGreen    = color.RGBA{0x00, 0xFF, 0x00, 0xFF}

	mapSize ixycord

	mapImg *ebiten.Image

	meltStart   *ebiten.Image
	meltBuf     *ebiten.Image
	meltOffsets [meltWidth]int
	screenSave  *ebiten.Image
	meltQuit    = false

	maxDof int

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
