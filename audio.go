package main

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var (
	mixer = &beep.Mixer{}
)

func init() {
	speaker.Init(48000, 100)
}

func getStream(file string) beep.StreamSeekCloser {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	streamer, _, err := mp3.Decode(f)
	if err != nil {
		panic(err)
	}
	return streamer
}

func playOst1() {
	stream := getStream("./assets/ost1.mp3")

	volume := &effects.Volume{Streamer: stream, Base: 2, Volume: -3.0}
	speaker.Play(volume)
	mixer.Add(stream)

	time.Sleep(time.Minute * 20)
	stream.Close()
}

func playWeaponSwing1() {
	stream := getStream("./assets/weapon_swing.mp3")
	mixer := &beep.Mixer{}

	speaker.Play(mixer)
	mixer.Add(stream)

	time.Sleep(time.Second * 2)
	stream.Close()
}

func playFootsteps(player *Player) {
	for {
		if player.isWalking {
			stream := getStream("./assets/footstep1.mp3")
			mixer := &beep.Mixer{}

			speaker.Play(mixer)
			mixer.Add(stream)

			time.Sleep(time.Millisecond * 500)
			stream.Close()
		}
	}
}
