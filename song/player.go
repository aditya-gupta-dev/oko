package song

import (
	"os"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type Player struct {
	streamer beep.StreamSeekCloser
	ctrl     *beep.Ctrl
	format   beep.Format
	mutex    sync.Mutex
	playing  bool
}

func NewPlayer() *Player {
	return &Player{}
}

func (p *Player) LoadFile(path string) error {
	p.Stop()
	p.mutex.Lock()
	defer p.mutex.Unlock()

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return err
	}
	ctrl := &beep.Ctrl{Streamer: beep.Loop(1, streamer), Paused: false}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(ctrl)

	p.streamer = streamer
	p.ctrl = ctrl
	p.format = format
	p.playing = true
	return nil
}

func (p *Player) Play() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.ctrl != nil {
		p.ctrl.Paused = false
		p.playing = true
	}
}

func (p *Player) Pause() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.ctrl != nil {
		p.ctrl.Paused = true
		p.playing = false
	}
}

func (p *Player) Toggle() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.ctrl != nil {
		p.ctrl.Paused = !p.ctrl.Paused
		p.playing = !p.ctrl.Paused
	}
}

func (p *Player) Seek(seconds int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.streamer != nil {
		samples := p.format.SampleRate.N(time.Duration(seconds) * time.Second)
		current := p.streamer.Position()
		p.streamer.Seek(current + samples)
	}
}

func (p *Player) Stop() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.streamer != nil {
		p.streamer.Close()
		p.streamer = nil
		p.ctrl = nil
		p.playing = false
	}
}

func (p *Player) IsPlaying() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.playing
}

func (p *Player) Position() (int, int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.streamer != nil {
		return p.streamer.Position(), p.streamer.Len()
	}
	return 0, 0
}

func (p *Player) Cleanup() {
	p.Stop()
}
