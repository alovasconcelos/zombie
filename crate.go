// Zombie
// crate.go
// Por André Luiz de Oliveira Vasconcelos (alovasconcelos@gmail.com)
// https://github.com/alovasconcelos/zombie
// 2020

package main

import "github.com/hajimehoshi/ebiten"

type crate struct {
	sprite *ebiten.Image
	x, y   float64
}

func (c *crate) draw(screen *ebiten.Image) {

	drawImage(screen, c.sprite, c.x, c.y)

}

// retorna true se o jogador esbarrar à esquerda de um obstaculo
func obstacleLeftEdge(position float64) bool {
	for i := 0; i < cap(obstacles); i++ {
		if int(position+76) == int(obstacles[i].x) {
			return true
		}
	}
	return false
}

// retorna true se o jogador esbarrar à direita de um obstaculo
func obstacleRightEdge(position float64) bool {
	for i := 0; i < cap(obstacles); i++ {
		if int(position+10) == int(obstacles[i].x) {
			return true
		}
	}
	return false
}

// retorna true se o jogador em cima de um obstaculo
func obstacleTop(x, y float64, forward bool) bool {
	for i := 0; i < cap(obstacles); i++ {
		if int(obstacles[i].y) == int(y+playerHeight) &&
			int(x+playerWidth/2) >= int(obstacles[i].x) &&
			int(x+playerWidth/2) <= int(obstacles[i].x+crateWidth) {
			return true
		}
	}
	return false
}
