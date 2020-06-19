// Zombie
// soundPlayer.go
// Por André Luiz de Oliveira Vasconcelos (alovasconcelos@gmail.com)
// https://github.com/alovasconcelos/zombie
// 2020

package main

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/mp3"
)

type soundPlayer struct {
	audioPlayer *audio.Player
}

func (s *soundPlayer) loadMP3(fileName string) {
	var err error

	// Initialize audio context.
	if audioContext == nil {
		audioContext, err = audio.NewContext(44100)
		if err != nil {
			log.Fatal(err)
		}
	}

	f, err := os.Open("assets/" + fileName)
	if err != nil {
		log.Fatal(err)
	}

	d, err := mp3.Decode(audioContext, f)

	if err != nil {
		log.Fatal(err)
	}

	s.audioPlayer, err = audio.NewPlayer(audioContext, d)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *soundPlayer) playMP3() {
	s.audioPlayer.Play()
}
