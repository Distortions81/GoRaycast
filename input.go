package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {
	return nil
}

func (g *Game) processInput(screen *ebiten.Image) error {

	/* Move/turn speed based on fps */
	var frameSpeed float64
	fps := ebiten.ActualFPS()
	if fps > 1 {
		frameSpeed = playerRotSpeed / fps
	} else {
		frameSpeed = playerRotSpeed / 60.0
	}

	/* Get pressed keys */
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	/* Process each key */
	for _, p := range g.keys {
		switch p {
		case ebiten.KeyF10:
			doMelt = -1
			meltQuit = true
		case ebiten.KeyW:
			playerPhysics.Position.x += playerPhysics.MovePos.x / playerForwardSpeedDiv
			playerPhysics.Position.y += playerPhysics.MovePos.y / playerForwardSpeedDiv
		case ebiten.KeyS:
			playerPhysics.Position.x -= playerPhysics.MovePos.x / playerForwardSpeedDiv
			playerPhysics.Position.y -= playerPhysics.MovePos.y / playerForwardSpeedDiv
		case ebiten.KeyD:
			playerPhysics.Rotation -= frameSpeed
			angleCalc() //Update player movepos
		case ebiten.KeyA:
			playerPhysics.Rotation += frameSpeed
			angleCalc() //Update player movepos
		}
	}
	return nil
}

func angleCalc() {
	playerPhysics.Rotation = fixRad(playerPhysics.Rotation)
	playerPhysics.MovePos.x = math.Cos(playerPhysics.Rotation)  // opposite
	playerPhysics.MovePos.y = -math.Sin(playerPhysics.Rotation) // adjacent
}
