package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {
	oldPlayerPos = playerPos

	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	pad := float64(1.0 / mapScale)
	for _, p := range g.keys {
		switch p {
		case ebiten.KeyS:
			if playerPos.y < float64(mapSize.y)-pad {
				playerPos.y += charMoveSpeed
			}

		case ebiten.KeyW:
			if playerPos.y > 0 {
				playerPos.y -= charMoveSpeed
			}

		case ebiten.KeyD:
			if playerPos.x < float64(mapSize.x)-pad {
				playerPos.x += charMoveSpeed
			}
		case ebiten.KeyA:
			if playerPos.x > 0 {
				playerPos.x -= charMoveSpeed
			}
		}
	}
	if oldPlayerPos != playerPos {
		playerImg.Fill(color.Transparent)
		playerImg.Set(int(playerPos.x*mapScale), int(playerPos.y*mapScale), cYellow)
	}
	return nil
}
