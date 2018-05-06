package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/faiface/beep"

	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

type decoder func(io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error)

var decoders = map[string]decoder{
	"wav":  wav.Decode,
	"mp3":  mp3.Decode,
	"flac": flac.Decode,
}

var (
	target            string
	sampleRateHz      int
	sampleRateQuality int
)

func main() {
	flag.StringVar(&target, "file", "", "File to play (wav, flac, or mp3)")
	flag.IntVar(&sampleRateHz, "resample-to", 48000, "Target Sample Rate Rate")
	flag.IntVar(&sampleRateQuality, "resample-quality", 4, "Resample Quality")
	flag.Parse()

	if target == "" {
		panic("File Required")
	}

	f, err := os.Open(target)
	if err != nil {
		panic(err)
	}

	ext := strings.Split(strings.ToLower(target), ".")
	d, ok := decoders[ext[len(ext)-1]]
	if !ok {
		panic(fmt.Sprintf("Unknown extension for file (%s): %s", ext, target))
	}

	s, format, err := d(f)
	if err != nil {
		panic(err)
	}

	done := make(chan struct{})

	speaker.UnderrunCallback(func() {
		fmt.Println("X-RUN: Not enough data in buffer. Get better hardware or increase buffer size")
	})

	fmt.Printf("Playing %s resampled to %dHz at quality level %d\n", target, sampleRateHz, sampleRateQuality)

	sr := beep.SampleRate(sampleRateHz)
	speaker.Init(sr, sr.N(time.Second/10))
	speaker.Play(
		beep.Seq(
			beep.Resample(sampleRateQuality, format.SampleRate, sr, s), beep.Callback(func() {
				close(done)
			}),
		),
	)

	<-done

	fmt.Println("Playback Complete")
}
