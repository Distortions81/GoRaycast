package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {

	playerRotSpeed = playerMoveSpeed / ebiten.ActualFPS()

	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	for _, p := range g.keys {
		switch p {
		case ebiten.KeyW:
			playerPhysics.Position.y += playerPhysics.MovePos.x
			playerPhysics.Position.x += playerPhysics.MovePos.y
		case ebiten.KeyS:
			playerPhysics.Position.y -= playerPhysics.MovePos.x
			playerPhysics.Position.x -= playerPhysics.MovePos.y
		case ebiten.KeyD:
			playerPhysics.Rotation += playerRotSpeed
			angleCalc()
		case ebiten.KeyA:
			playerPhysics.Rotation -= playerRotSpeed
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
	playerPhysics.MovePos.x = math.Cos(playerPhysics.Rotation)
	playerPhysics.MovePos.y = -math.Sin(playerPhysics.Rotation)
}
