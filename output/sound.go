package output

import (
	"time"

	"github.com/veandco/go-sdl2/mix"
)

// PlaySound plays a sound file with specified volume and loop settings.
func PlaySound(filename string, volume int, loop bool) {
	go func() {
		chunk, err := mix.LoadWAV(filename)
		if err != nil {
			Printf("Failed to load sound: %s\n", err)
			return
		}

		chunk.Volume(volume)

		if loop {
			for {
				chunk.Play(-1, 0)
				time.Sleep(time.Duration(chunk.LengthInMs()) * time.Millisecond)
			}
		} else {
			chunk.Play(-1, 0)
		}
	}()
}
