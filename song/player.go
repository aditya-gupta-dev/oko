package song

import (
	"os"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type Player struct {
	streamer      beep.StreamSeekCloser
	ctrl          *beep.Ctrl
	volume        *effects.Volume
	format        beep.Format
	mutex         sync.Mutex
	playing       bool
	loaded        bool
	looping       bool
	volumePercent int
}

func NewPlayer() *Player {
	return &Player{
		volumePercent: 70,
	}
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
	p.looping = false
	p.volume = &effects.Volume{
		Streamer: p.wrapStreamer(streamer),
		Base:     2,
		Volume:   p.volumeLevelFromPercent(),
		Silent:   p.volumePercent == 0,
	}
	ctrl := &beep.Ctrl{Streamer: p.volume, Paused: false}

	if !p.loaded {
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		p.loaded = true
	}
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
		p.volume = nil
		p.playing = false
		p.looping = false
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

func (p *Player) samplesToDuration(samples int) time.Duration {
	if p.format.SampleRate <= 0 {
		return 0
	}
	seconds := float64(samples) / float64(p.format.SampleRate)
	return time.Duration(seconds * float64(time.Second))
}

func (p *Player) PositionDuration() (time.Duration, time.Duration) {
	currentSamples, totalSamples := p.Position()
	currentDuration := p.samplesToDuration(currentSamples)
	totalDuration := p.samplesToDuration(totalSamples)
	return currentDuration, totalDuration
}

func (p *Player) Cleanup() {
	p.Stop()
}

func (p *Player) ToggleLoop() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.ctrl == nil || p.streamer == nil {
		return p.looping
	}

	p.looping = !p.looping

	speaker.Lock()
	if p.volume != nil {
		p.volume.Streamer = p.wrapStreamer(p.streamer)
		p.ctrl.Streamer = p.volume
	}
	speaker.Unlock()

	return p.looping
}

func (p *Player) IsLooping() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.looping
}

func (p *Player) HasTrackLoaded() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.ctrl != nil && p.streamer != nil
}

func (p *Player) wrapStreamer(streamer beep.StreamSeekCloser) beep.Streamer {
	if p.looping {
		return beep.Loop(-1, streamer)
	}

	return beep.Loop(1, streamer)
}

func (p *Player) IncreaseVolume() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.volumePercent < 100 {
		p.volumePercent += 5
		if p.volumePercent > 100 {
			p.volumePercent = 100
		}
	}

	p.applyVolume()

	return p.volumePercent
}

func (p *Player) DecreaseVolume() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.volumePercent > 0 {
		p.volumePercent -= 5
		if p.volumePercent < 0 {
			p.volumePercent = 0
		}
	}

	p.applyVolume()

	return p.volumePercent
}

func (p *Player) GetVolumePercent() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.volumePercent
}

func (p *Player) applyVolume() {
	if p.volume == nil {
		return
	}

	speaker.Lock()
	p.volume.Volume = p.volumeLevelFromPercent()
	p.volume.Silent = p.volumePercent == 0
	speaker.Unlock()
}

func (p *Player) volumeLevelFromPercent() float64 {
	if p.volumePercent <= 0 {
		return -8
	}

	return (float64(p.volumePercent) / 100.0 * 8.0) - 8.0
}
