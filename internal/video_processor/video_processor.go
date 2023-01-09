package video_processor

import (
	"fmt"
	"gocv.io/x/gocv"
	"time"
)

type ProcessingUnit struct {
	imgCh    chan inputStruct
	finishCh chan struct{}
	fname    string
	writer   *gocv.VideoWriter
	active   bool
}

func (pu *ProcessingUnit) Active() bool {
	return pu.active
}

type inputStruct struct {
	img    gocv.Mat
	numObj int
}

func NewProcessingUnit() (*ProcessingUnit, error) {
	fname := fmt.Sprintf("%s.mp4", time.Now().Format("20060102150405"))
	writer, err := gocv.VideoWriterFile(fname, "H264", 25, 1280, 720, true)
	if err != nil {
		return nil, err
	}
	return &ProcessingUnit{
		imgCh:    make(chan inputStruct),
		finishCh: make(chan struct{}),
		fname:    fname,
		writer:   writer,
		active:   true,
	}, nil
}

func (pu *ProcessingUnit) Send(img gocv.Mat, numObj int) {
	pu.imgCh <- inputStruct{
		img,
		numObj,
	}
}

func (pu *ProcessingUnit) Stop() {
	pu.finishCh <- struct{}{}
}

func (pu *ProcessingUnit) Processing() {
	select {
	case in := <-pu.imgCh:
		if in.numObj == 0 {
			pu.Stop()
		}
		pu.writer.Write(in.img)
		break
	case <-pu.finishCh:
		close(pu.imgCh)
		close(pu.finishCh)
		pu.writer.Close()
		pu.active = false
		return
	}
}
