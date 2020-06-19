// Zombie
// skeleton.go
// Por André Luiz de Oliveira Vasconcelos (alovasconcelos@gmail.com)
// https://github.com/alovasconcelos/zombie
// 2020

package main

import (
	"time"

	"github.com/hajimehoshi/ebiten"
)

type skeleton struct {
	sprite              *ebiten.Image
	x                   float64
	y                   float64
	speed               float64
	killedTheZombie     bool
	active              bool
	skeletonUpdateDelay time.Duration
	lastUpdate          time.Time
}

func (s *skeleton) update(p *player) {
	if time.Since(s.lastUpdate) >= s.skeletonUpdateDelay {
		if s.x <= 0 || s.killedTheZombie {
			s.x = screenWidth - skeletonWidth

			s.killedTheZombie = false
		}
		s.x -= s.speed
		if s.x <= p.x && s.x >= p.x-playerWidth/2 && p.y == screenHeight-playerHeight {
			p.points++
		}
		s.lastUpdate = time.Now()
	}
}
