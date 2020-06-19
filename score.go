// Zombie
// score.go
// Por André Luiz de Oliveira Vasconcelos (alovasconcelos@gmail.com)
// https://github.com/alovasconcelos/zombie
// 2020

package main

import "github.com/hajimehoshi/ebiten"

type score struct {
	sprite   *ebiten.Image
	lives    int
	gameOver bool
}

func (s *score) update(screen *ebiten.Image) {
	for i := 0; i < s.lives; i++ {
		drawImage(screen, s.sprite, float64(i*40), 55)
	}

}

func (s *score) hit() {
	if s.lives > 0 {
		s.gameOver = false
		s.lives--
	} else {
		// game over
		s.gameOver = true
	}

}
