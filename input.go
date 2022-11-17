package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {

	playerRotSpeed = 10.0 / ebiten.ActualFPS()

	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	for _, p := range g.keys {
		switch p {
		case ebiten.KeyW:
			playerPos.y += playerPosR.x
			playerPos.x += playerPosR.y
		case ebiten.KeyS:
			playerPos.y -= playerPosR.x
			playerPos.x -= playerPosR.y
		case ebiten.KeyD:
			playerRot += playerRotSpeed
			if playerRot > twoPi {
				playerRot -= twoPi
			}
			playerPosR.x = math.Cos(playerRot)
			playerPosR.y = -math.Sin(playerRot)
		case ebiten.KeyA:
			playerRot -= playerRotSpeed
			if playerRot > twoPi {
				playerRot -= twoPi
			}
			playerPosR.x = math.Cos(playerRot)
			playerPosR.y = -math.Sin(playerRot)
		}
	}
	return nil
}
