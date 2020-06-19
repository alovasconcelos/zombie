// Zombie
// player.go
// Por André Luiz de Oliveira Vasconcelos (alovasconcelos@gmail.com)
// https://github.com/alovasconcelos/zombie
// 2020

package main

import (
	"image"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	stateIdle = "Idle"
	stateWalk = "Walk"
	stateDead = "Dead"
	stateHurt = "Hurt"
	stateDone = "Done"

	framesIdle = 12
	framesWalk = 8
	framesDead = 8
)

type player struct {
	idleSprites  [framesIdle]*ebiten.Image
	idleRSprites [framesIdle]*ebiten.Image
	walkSprites  [framesWalk]*ebiten.Image
	walkRSprites [framesWalk]*ebiten.Image
	deadSprites  [framesDead]*ebiten.Image
	deadRSprites [framesDead]*ebiten.Image

	active bool

	x, y  float64
	speed float64

	state string
	frame int

	forward bool
	jumping bool
	falling bool
	paused  bool
	help    bool
	exit    bool
	finish  bool

	sk [3]skeleton
	sc score

	zombieGroan soundPlayer
	gameOver    soundPlayer
	ouch        soundPlayer
	jumpSound   soundPlayer

	lastTimeUpdate time.Time
	timeLeft       int

	points int
}

// carrega os sprites para a animação do zumbi
func (p *player) loadSprites() {
	for i := 0; i < framesIdle; i++ {
		p.idleSprites[i], _ = initializeImage("assets/Idle" + strconv.Itoa(i+1) + ".png")
		p.idleRSprites[i], _ = initializeImage("assets/IdleR" + strconv.Itoa(i+1) + ".png")
	}
	for i := 0; i < framesWalk; i++ {
		p.walkSprites[i], _ = initializeImage("assets/Walk" + strconv.Itoa(i+1) + ".png")
		p.walkRSprites[i], _ = initializeImage("assets/WalkR" + strconv.Itoa(i+1) + ".png")
	}
	for i := 0; i < framesDead; i++ {
		p.deadSprites[i], _ = initializeImage("assets/Die" + strconv.Itoa(i+1) + ".png")
		p.deadRSprites[i], _ = initializeImage("assets/DieR" + strconv.Itoa(i+1) + ".png")
	}

}

func (p *player) level1(screen *ebiten.Image) {

	if time.Since(p.lastTimeUpdate) >= time.Second {
		p.lastTimeUpdate = time.Now()
		p.timeLeft--
	}

	drawImage(screen, bg, 0, 0)
	drawImage(screen, logo, screenWidth-240, 0)
	drawImage(screen, ubLogo, screenWidth-100, 100)
	drawImage(screen, f1, 0, 10)
	drawImage(screen, arrowSign, 80, screenHeight-86)
	cx := [5]int{200, 400, 550, 700, 850}
	cr.y = screenHeight - crateHeight

	for i := 0; i < 5; i++ {
		cr.x = float64(cx[i])
		cr.draw(screen)
		obstacles[i] = cr
	}
	p.update(screen)
}

// atualização da tela do jogo
func (p *player) update(screen *ebiten.Image) {

	// Nosso herói conseguiu atravessar toda a fase
	if p.state == stateDone {
		// mostra a imagem de fim de jogo
		drawImage(screen, bigTombstone, screenWidth/2-bigTombstoneWidth/2, 10)
		drawImage(screen, doneImage, screenWidth/2-doneWidth/2+5, screenHeight/3-doneHeight/2)
		ebitenutil.DebugPrintAt(screen, "Seu placar foi:", 480, screenHeight/2+40)
		showNumber(screen, p.points, screenWidth/2-48, screenHeight/2+60)
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			p.active = false
		}

		return
	}

	p.draw(screen)

	if ebiten.IsKeyPressed(ebiten.KeyF6) {
		p.paused = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyF7) {
		p.paused = false
	}
	if p.paused || p.help || p.exit {
		return
	}
	if ebiten.IsKeyPressed(ebiten.KeyF1) {
		helpActivatedAt = time.Now()
		p.help = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		p.exit = true
	}
	if p.state != stateDead {
		if p.jumping {
			if p.y < screenHeight-(playerHeight*1.5) {
				p.jumping = false
				p.falling = true
			} else {
				p.y -= p.speed * 4
			}
		}

		if p.falling {
			if p.y == screenHeight-playerHeight {
				p.falling = false
			} else if !obstacleTop(p.x, p.y, p.forward) {
				p.y += p.speed * 4
			}
		}

		// Up - pular
		if ebiten.IsKeyPressed(ebiten.KeyUp) && !p.jumping && !p.falling {
			p.jump()
		}
		// Left - voltar para o lado esquerdo
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			if p.x > 10 && (!obstacleRightEdge(p.x) || p.jumping || p.falling) {
				p.x -= p.speed
				p.goBack()
			}
		}
		// Right - avançar para o lado direito
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			if p.x+playerWidth < screenWidth && (!obstacleLeftEdge(p.x) || p.jumping || p.falling) {
				p.x += p.speed
				p.goForward()
			}
		}

		if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyRight) {
			// se não estava andando, passa para o modo andando
			if p.state != stateWalk {
				p.frame = 0
			}
			p.state = stateWalk
		} else if p.state == stateWalk {
			// se estava andando passa para o modo parado
			p.state = stateIdle
		}
	}
	if time.Since(lastFrameUpdate) >= frameUpdateDelay && (!p.jumping && !p.falling || p.state == stateDead) {
		if p.state != stateDead && p.state != stateHurt && p.state != stateDone {
			if ebiten.IsKeyPressed(ebiten.KeyF2) {
				muteSoundtrack = true
				soundtrack.audioPlayer.Pause()
			}
			if ebiten.IsKeyPressed(ebiten.KeyF3) {
				muteSoundtrack = false
				soundtrack.audioPlayer.Play()
			}

			if ebiten.IsKeyPressed(ebiten.KeyF4) {
				muteGroan = true
				p.zombieGroan.audioPlayer.Pause()
			}
			if ebiten.IsKeyPressed(ebiten.KeyF5) {
				muteGroan = false
				p.zombieGroan.audioPlayer.Play()
			}
		} else {
			if ebiten.IsKeyPressed(ebiten.KeyEnter) {
				p.active = false
			}
		}

		p.frame++

		switch {
		case p.state == stateIdle && p.frame == framesIdle:
			p.frame = 0
		case p.state == stateWalk && p.frame == framesWalk:
			p.frame = 0
		case p.state == stateHurt && p.frame == framesIdle/3:
			p.frame = 0
			p.state = stateIdle
		case p.state == stateDead && p.frame == framesDead:
			// fim de jogo
			p.frame = framesDead - 1
		}
		lastFrameUpdate = time.Now()
	}

	// 'mutar' o zumbi
	if p.frame > 7 && !p.zombieGroan.audioPlayer.IsPlaying() && !muteGroan {
		p.zombieGroan.audioPlayer.Rewind()
		p.zombieGroan.playMP3()
	}

}

// desenhar os elementos do jogo
func (p *player) draw(screen *ebiten.Image) {
	var zombieImage *ebiten.Image

	zombieImage = p.idleRSprites[0]
	switch p.state {
	case stateIdle:
		if p.forward {
			zombieImage = p.idleSprites[p.frame]
		} else {
			zombieImage = p.idleRSprites[p.frame]
		}
	case stateWalk:
		if p.forward {
			zombieImage = p.walkSprites[p.frame]
		} else {
			zombieImage = p.walkRSprites[p.frame]
		}
	case stateDead:
		if p.forward {
			zombieImage = p.deadSprites[p.frame]
		} else {
			zombieImage = p.deadRSprites[p.frame]
		}
	case stateHurt:
		if p.forward {
			zombieImage = p.idleSprites[0]
		} else {
			zombieImage = p.idleRSprites[0]
		}
	}
	if p.state != stateDead && p.timeLeft == 0 {
		drawImage(screen, zombieImage, p.x, p.y)
		p.die()
		return
	}
	if p.paused {
		return
	}
	if p.help {
		drawImage(screen, helpImage, screenWidth/2-helpWidth/2, 10)
		if time.Since(helpActivatedAt) >= time.Second*10 ||
			ebiten.IsKeyPressed(ebiten.KeyEnter) {
			p.help = false
		}
		return
	}
	if p.exit {
		drawImage(screen, confirm, screenWidth/2-confirmWidth/2, screenHeight/3-confirmHeight/2)
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			p.finish = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyN) {
			p.exit = false
		}
		return
	}
	zombieHit := false
	// atualiza posição dos esqueletos
	for x := 0; x < cap(p.sk); x++ {
		if x > 0 {
			if !p.sk[x].active && p.sk[x-1].x < screenWidth-skeletonWidth-300 {
				p.sk[x].active = true
			}
		}

		if !p.sk[x].active {
			continue
		}

		p.sk[x].update(p)
		// desenha o esqueleto
		drawImage(screen, p.sk[x].sprite, p.sk[x].x, p.sk[x].y)
		// verifica se o esqueleto atingiu o zumbi
		if int(p.sk[x].x) >= int(p.x) && int(p.sk[x].x) <= int(p.x+playerWidth/2) && !p.sk[x].killedTheZombie && p.y == screenHeight-playerHeight {
			zombieHit = true
			p.sk[x].killedTheZombie = true
		}
	}

	// desenha o zumbi
	if p.state != stateHurt || p.frame%2 == 0 {
		drawImage(screen, zombieImage, p.x, p.y)
	}

	// o zumbi não está morto (não diga...) e foi atingido
	if p.state != stateDead && zombieHit {
		// se o zumbi não estiver "mudo", ele grita
		if !muteGroan {
			p.ouch.audioPlayer.Rewind()
			p.ouch.playMP3()
		}
		p.sc.hit()
		p.frame = 0
		p.state = stateHurt
		if p.sc.gameOver {
			// fim do jogo
			p.die()
			return
		}
	}

	// o zumbi está morto e está no final da animação (último frame)
	if p.state == stateDead && p.frame == framesDead-1 {
		// mostra a imagem de fim de jogo
		drawImage(screen, bigTombstone, screenWidth/2-bigTombstoneWidth/2, 10)
		drawImage(screen, gameOverImage, screenWidth/2-gameOverWidth/2+5, screenHeight/3-gameOverHeight/2)
		ebitenutil.DebugPrintAt(screen, "Seu placar foi:", 480, screenHeight/2+40)
		showNumber(screen, p.points, screenWidth/2-48, screenHeight/2+60)
	} else {
		// mostra o tempo restante
		if p.timeLeft >= 0 {
			p.showTimeLeft(screen)
		}
	}

	// o zumbi conseguiu atravessar toda a fase
	if p.state != stateDead && p.x >= screenWidth-playerWidth-50 {
		p.state = stateDone
	}

	// atualiza o placar
	p.sc.update(screen)
}

// avançar
func (p *player) goForward() {
	p.forward = true
}

// voltar
func (p *player) goBack() {
	p.forward = false
}

// pular
func (p *player) jump() {
	p.jumpSound.audioPlayer.Rewind()
	p.jumpSound.playMP3()
	p.jumping = true
	p.falling = false
}

// ocioso
func (p *player) idle() {
	p.state = stateIdle
	p.frame = 0
}

// morre (como assim?! - ele é um zumbi!!!)
func (p *player) die() {
	soundtrack.audioPlayer.Pause()
	zombiePlayer.zombieGroan.audioPlayer.Pause()
	zombiePlayer.gameOver.audioPlayer.Rewind()
	zombiePlayer.gameOver.playMP3()
	p.state = stateDead
	p.frame = 0
}

// mostra o tempo restante
func (p *player) showTimeLeft(screen *ebiten.Image) {
	drawImage(screen, timeStr, 15, 120)
	showNumber(screen, p.timeLeft, 10, 140)

	drawImage(screen, pointsStr, 15, 220)
	showNumber(screen, p.points, 10, 240)
}

// mostra número
func showNumber(screen *ebiten.Image, number, x, y int) {
	strNumber := strconv.Itoa(number)
	for i := 0; i < len(strNumber); i++ {
		d, _ := strconv.Atoi(string(strNumber[i]))
		drawImage(screen, numbers.SubImage(image.Rect(d*57, 0, d*57+57, 66)).(*ebiten.Image),
			float64(x+53*i), float64(y))
	}

}
