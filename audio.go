package main

import (
	"fmt"
	"io"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/jonathangjertsen/toyboy/model"
)

type Audio struct {
	In     chan []model.AudioSample
	buf    []uint8
	offset int
	player *oto.Player
}

func NewAudio() *Audio {
	audio := &Audio{}
	audio.In = make(chan []model.AudioSample, 2*44100/1024)
	opts := &oto.NewContextOptions{}
	opts.SampleRate = 44100
	opts.ChannelCount = 1
	opts.BufferSize = time.Millisecond
	opts.Format = oto.FormatSignedInt16LE
	otoCtx, readyChan, err := oto.NewContext(opts)
	if err != nil {
		panic("oto.NewContext failed with " + err.Error())
	}
	<-readyChan
	audio.player = otoCtx.NewPlayer(audio)
	//go audio.player.Play()
	return audio
}

func (audio *Audio) Read(p []byte) (int, error) {
	nTotal := len(p)
	nRemaining := nTotal
	nFilled := 0
	for nRemaining > 0 {
		if len(audio.buf) == 0 {
			select {
			case data, ok := <-audio.In:
				if !ok {
					audio.player.Close()
					return nFilled, io.EOF
				}
				for s := range data {
					s <<= 10
					audio.buf = append(audio.buf, uint8(s), uint8(s>>8))
				}
			default:
				fmt.Printf("EMPTY\n")
				return nFilled, nil
			}
		}
		nBuf := len(audio.buf)
		if nRemaining >= nBuf {
			copy(p[nFilled:nFilled+nBuf], audio.buf)
			audio.buf = make([]uint8, 0, 1024)
			nRemaining -= nBuf
			nFilled += nBuf
		} else {
			copy(p[nFilled:nFilled+nRemaining], audio.buf[:nRemaining])
			audio.buf = audio.buf[nRemaining:]
			nFilled += nRemaining
			nRemaining = 0
		}
	}
	return nFilled, nil
}
