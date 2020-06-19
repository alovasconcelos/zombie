// Zombie
// main.go
// Por André Luiz de Oliveira Vasconcelos (alovasconcelos@gmail.com)
// https://github.com/alovasconcelos/zombie
// 2020

package main

import (
	"image"
	_ "image/png"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth        = 1024
	screenHeight       = 585
	screenScale        = 1
	playerHeight       = 87
	playerWidth        = 100
	crateHeight        = 30
	crateWidth         = 30
	skeletonWidth      = 30
	skeletonHeight     = 15
	gameOverWidth      = 140
	gameOverHeight     = 226
	doneWidth          = 273
	doneHeight         = 61
	helpWidth          = 488
	helpHeight         = 519
	confirmWidth       = 562
	confirmHeight      = 138
	bigTombstoneWidth  = 398
	bigTombstoneHeight = 500
	appTitle           = "Zombie"
	logoWidth          = 234
	frameUpdateDelay   = time.Millisecond * 100
)

var (
	err           error
	bg            *ebiten.Image
	logo          *ebiten.Image
	zombie        *ebiten.Image
	arrowSign     *ebiten.Image
	gameOverImage *ebiten.Image
	doneImage     *ebiten.Image
	helpImage     *ebiten.Image
	f1            *ebiten.Image
	confirm       *ebiten.Image
	bigTombstone  *ebiten.Image
	numbers       *ebiten.Image
	pointsStr     *ebiten.Image
	timeStr       *ebiten.Image
	ubLogo        *ebiten.Image
	cr            crate
	obstacles     [5]crate

	lastFrameUpdate time.Time
	helpActivatedAt time.Time
	level           int
	soundtrack      soundPlayer
	zombiePlayer    player
	audioContext    *audio.Context

	scor score

	muteSoundtrack bool
	muteGroan      bool

	tempoJogo int
)

func start() {

	// nível
	level = 1

	// logo
	logo, _ = initializeImage("assets/logo.png")

	// sinal
	arrowSign, _ = initializeImage("assets/ArrowSign.png")

	// game over
	gameOverImage, _ = initializeImage("assets/gameOver.png")

	// done
	doneImage, _ = initializeImage("assets/done.png")

	// help
	helpImage, _ = initializeImage("assets/help.png")

	// f1
	f1, _ = initializeImage("assets/f1.png")

	// confirm
	confirm, _ = initializeImage("assets/confirm.png")

	// big tombstone
	bigTombstone, _ = initializeImage("assets/tombstone.png")

	// background
	bg, _ = initializeImage("assets/bg.png")

	// numbers
	numbers, _ = initializeImage("assets/numbers.png")

	// points
	pointsStr, _ = initializeImage("assets/points.png")

	// time
	timeStr, _ = initializeImage("assets/time.png")

	// Useless Bytes logo
	ubLogo, _ = initializeImage("assets/ublogo.png")

	lastFrameUpdate = time.Now()

	tempoJogo = 35

	scor = score{}
	scor.sprite, _ = initializeImage("assets/Head.png")
	scor.lives = 3

	zombiePlayer = player{
		speed:   .5,
		x:       10,
		y:       screenHeight - playerHeight,
		state:   stateIdle,
		forward: true,
		jumping: false,
		falling: false,
		sc:      scor,
		active:  true,
		paused:  false,
		help:    false,
		exit:    false,
		finish:  false,
	}
	zombiePlayer.zombieGroan = soundPlayer{}
	zombiePlayer.zombieGroan.loadMP3("zombie-groan.mp3")
	zombiePlayer.zombieGroan.audioPlayer.SetVolume(.4)
	zombiePlayer.zombieGroan.playMP3()
	zombiePlayer.gameOver = soundPlayer{}
	zombiePlayer.gameOver.loadMP3("game-over.mp3")
	zombiePlayer.ouch = soundPlayer{}
	zombiePlayer.ouch.loadMP3("ouch.mp3")
	zombiePlayer.jumpSound = soundPlayer{}
	zombiePlayer.jumpSound.loadMP3("jump.mp3")
	zombiePlayer.jumpSound.audioPlayer.SetVolume(.1)
	zombiePlayer.ouch.audioPlayer.SetVolume(.4)
	zombiePlayer.lastTimeUpdate = time.Now()
	zombiePlayer.timeLeft = tempoJogo
	zombiePlayer.points = 0
	zombiePlayer.idle()

	for x := 0; x < 3; x++ {
		zombiePlayer.sk[x] = skeleton{
			x:                   screenWidth - skeletonWidth,
			y:                   screenHeight - skeletonHeight,
			speed:               2,
			killedTheZombie:     false,
			lastUpdate:          time.Now(),
			skeletonUpdateDelay: time.Millisecond * 10,
		}
		if x == 0 {
			zombiePlayer.sk[x].active = true
		} else {
			zombiePlayer.sk[x].active = false
		}
		zombiePlayer.sk[x].sprite, _ = initializeImage("assets/Skeleton.png")
	}
	muteSoundtrack = false
	muteGroan = false

	// carrega sprites do zumbi
	zombiePlayer.loadSprites()

	// carrega sprite do engradado
	cr = crate{}
	cr.sprite, _ = initializeImage("assets/Crate.png")

	soundtrack = soundPlayer{}

	soundtrack.loadMP3("zombie.mp3")
	soundtrack.audioPlayer.SetVolume(.3)
	soundtrack.playMP3()

}

func init() {
	ebiten.SetWindowDecorated(false)
	start()
}

func restart() {
	// nível
	level = 1

	lastFrameUpdate = time.Now()

	zombiePlayer.sc.lives = 3
	zombiePlayer.sc.gameOver = false
	zombiePlayer.active = true
	zombiePlayer.state = stateIdle
	zombiePlayer.frame = 0
	zombiePlayer.x = 10
	zombiePlayer.forward = true
	zombiePlayer.paused = false
	zombiePlayer.help = false
	zombiePlayer.exit = false
	zombiePlayer.finish = false
	zombiePlayer.lastTimeUpdate = time.Now()
	zombiePlayer.timeLeft = tempoJogo
	zombiePlayer.points = 0

	zombiePlayer.sk[0].x = screenWidth - skeletonWidth
	zombiePlayer.sk[0].y = screenHeight - skeletonHeight
	zombiePlayer.sk[0].speed = 2
	zombiePlayer.sk[0].killedTheZombie = false
	zombiePlayer.sk[0].lastUpdate = time.Now()
	zombiePlayer.sk[0].active = true

	zombiePlayer.sk[1].x = screenWidth - skeletonWidth
	zombiePlayer.sk[1].y = screenHeight - skeletonHeight
	zombiePlayer.sk[1].speed = 2
	zombiePlayer.sk[1].killedTheZombie = false
	zombiePlayer.sk[1].lastUpdate = time.Now()
	zombiePlayer.sk[1].active = false

	zombiePlayer.sk[2].x = screenWidth - skeletonWidth
	zombiePlayer.sk[2].y = screenHeight - skeletonHeight
	zombiePlayer.sk[2].speed = 2
	zombiePlayer.sk[2].killedTheZombie = false
	zombiePlayer.sk[2].lastUpdate = time.Now()
	zombiePlayer.sk[2].active = false

	zombiePlayer.y = screenHeight - playerHeight

	if !muteSoundtrack {
		soundtrack.audioPlayer.Rewind()
		soundtrack.playMP3()
	}
	zombiePlayer.state = stateIdle
}

func initializeImage(fileName string) (*ebiten.Image, image.Image) {
	ebImage, image, err := ebitenutil.NewImageFromFile(fileName, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	return ebImage, image
}

func drawImage(screen *ebiten.Image, image *ebiten.Image, x, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	screen.DrawImage(image, op)

}

func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	switch level {
	case 1:
		zombiePlayer.level1(screen)
	}

	if !zombiePlayer.active {
		restart()
	}

	if zombiePlayer.finish {
		os.Exit(2)
	}

	// "mutar" a trilha do jogo
	if !zombiePlayer.sc.gameOver && !muteSoundtrack && !soundtrack.audioPlayer.IsPlaying() {
		soundtrack.audioPlayer.Rewind()
		soundtrack.playMP3()
	}

	return nil
}

func main() {
	if err := ebiten.Run(update,
		screenWidth,
		screenHeight,
		screenScale,
		appTitle); err != nil {
		log.Fatal(err)
	}
}
