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
			playerData.Position.x += playerData.MovePos.x / playerForwardSpeedDiv
			playerData.Position.y += playerData.MovePos.y / playerForwardSpeedDiv
		case ebiten.KeyS:
			playerData.Position.x -= playerData.MovePos.x / playerForwardSpeedDiv
			playerData.Position.y -= playerData.MovePos.y / playerForwardSpeedDiv
		case ebiten.KeyD:
			playerData.Rotation -= frameSpeed
			angleCalc() //Update player movepos
		case ebiten.KeyA:
			playerData.Rotation += frameSpeed
			angleCalc() //Update player movepos
		}
	}
	return nil
}

func angleCalc() {
	playerData.Rotation = fixRad(playerData.Rotation)
	playerData.MovePos.x = math.Cos(playerData.Rotation)  // opposite
	playerData.MovePos.y = -math.Sin(playerData.Rotation) // adjacent
}
