package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {

	frameSpeed := playerRotSpeed / ebiten.ActualFPS()

	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	for _, p := range g.keys {
		switch p {
		case ebiten.KeyW:
			playerPhysics.Position.x += playerPhysics.MovePos.x / playerForwardSpeedDiv
			playerPhysics.Position.y += playerPhysics.MovePos.y / playerForwardSpeedDiv
		case ebiten.KeyS:
			playerPhysics.Position.x -= playerPhysics.MovePos.x / playerForwardSpeedDiv
			playerPhysics.Position.y -= playerPhysics.MovePos.y / playerForwardSpeedDiv
		case ebiten.KeyD:
			playerPhysics.Rotation -= frameSpeed
			angleCalc()
		case ebiten.KeyA:
			playerPhysics.Rotation += frameSpeed
			angleCalc()
		}
	}
	return nil
}

func angleCalc() {
	if playerPhysics.Rotation > twoPi {
		playerPhysics.Rotation -= twoPi
	} else if playerPhysics.Rotation < 0 {
		playerPhysics.Rotation += twoPi
	}
	playerPhysics.MovePos.x = math.Cos(playerPhysics.Rotation)  // opposite
	playerPhysics.MovePos.y = -math.Sin(playerPhysics.Rotation) // adjacent
}
