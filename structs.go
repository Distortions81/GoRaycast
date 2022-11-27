package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	maxDist         = 1000000.0 //Used to signify no wall found
	renderFov       = 60        //Degrees
	screenWidth     = 1280
	screenHeight    = 720
	meltWidth       = screenWidth
	meltHeight      = screenHeight
	wallHeightRatio = 1.0
	screenMag       = 1 //Maginify screen, mosaic
	mapScale        = 1 //Units per map pixel

	/*
	 * Offset distance, so less than one unit
	 * of space does not over expose
	 */
	distanceOffset = shadowDistance / 4 //added to distance
	shadowDistance = 16                 //Shadow is divide by this, sets how far we can see
	FalloffRatio   = 16                 //Square root multiplied by this
	shadowBase     = 2                  //The shadow is divided by this
	shadowExp      = 2.4                //Lighting exponent value
	shadowClip     = 1.2                //Don't go brighter than this value

	/* Shade horizontal walls a bit, faux shading */
	normalShading = 1.0
	dirShading    = 0.85 //1.0 no shading, 0.85 darkens by 15%

	/* Player rotate/move speed */
	playerRotSpeed        = 2
	playerForwardSpeedDiv = 10

	/* Minimap */
	miniScale  = screenWidth / 120
	miniRayMod = screenWidth / 320

	/* Commonly used radians values */
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
	doMelt       int     //-1 to start timer, otherwise number of frames remaining
	renderFovRad float64 //FoV in radians
	halfFovRad   float64 //Half fov, to setup
	radPerRay    float64 //Radians to add per ray

	/* Some predfined colors */
	cDarkGray = color.RGBA{0x20, 0x20, 0x20, 0xFF}
	cYellow   = color.RGBA{0xFF, 0xAA, 0x00, 0xFF}
	cRed      = color.RGBA{0xFF, 0x00, 0x00, 0xFF}
	cGreen    = color.RGBA{0x00, 0xFF, 0x00, 0xFF}
	cSmoked   = color.RGBA{0x00, 0x00, 0x00, 0xFF}

	cRay = color.RGBA{0x60, 0x50, 0x30, 0xFF}
	cBG  = color.RGBA{0x20, 0x20, 0x20, 0xFF}

	/* Map size, and source image */
	mapSize  ixycord
	wallSize ixycord
	mapImg   *ebiten.Image
	titleImg *ebiten.Image
	wallImg  *ebiten.Image
	rayImg   *ebiten.Image
	miniMap  *ebiten.Image

	/* Screen melt buffers and offsets */
	meltStart   *ebiten.Image  //Converted starting image
	meltBuf     *ebiten.Image  //Melt effect output
	meltOffsets [meltWidth]int //Pixel offsets for melt effect

	screenSave *ebiten.Image //Screen capture
	meltQuit   = false       //If true, melt is for exiting game
	meltDelay  = 120         //How long to draw first frame for
	/* Screen melt params */
	meltSpeed = (screenHeight / 200) * 4
	maxDof    int //Max depth

	playerPhysics pPhysics //Player pos/rot/movepos
)

/* Player info */
type pPhysics struct {
	Position xycord
	Rotation float64
	MovePos  xycord
}

/* Keys pressed */
type Game struct {
	keys []ebiten.Key
}

/* x/y float64 */
type xycord struct {
	x float64
	y float64
}

/* x/y integer */
type ixycord struct {
	x int
	y int
}
