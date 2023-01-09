package internal

import (
	"context"
	"fmt"
	"gocv.io/x/gocv"
	"time"
)

func Subrecord(ctx context.Context, frames []gocv.Mat, sendCh chan string) error {
	fname := fmt.Sprintf("%s.avi", time.Now().Format("20060102150405"))
	writer, err := gocv.VideoWriterFile(fname, "H264", 25, frames[0].Cols(), frames[0].Rows(), true)
	if err != nil {
		return err
	}

	for _, frame := range frames {
		writer.Write(frame)
	}

	writer.Close()
	sendCh <- fname

	return nil
}
