package video_processor

import (
	"context"
	"errors"
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"log"
	"time"
)

type Processor struct {
	config ProcessorConfig
}

type ProcessorConfig struct {
	MinimumArea float64
}

func NewProcessor(conf ProcessorConfig) *Processor {
	return &Processor{
		config: conf,
	}
}

var (
	ErrCantReadImg  = errors.New("cant read img")
	ErrDeviceClosed = errors.New("device closed")
)

func subrecord(ctx context.Context, frames []gocv.Mat, sendCh chan string) error {
	fname := fmt.Sprintf("%s.mp4", time.Now().Format("20060102150405"))
	//writer, err := gocv.VideoWriterFile(fname, "H264", 25, frames[0].Cols(), frames[0].Rows(), true)
	writer, err := gocv.VideoWriterFile(fname, "H264", 25, 480, 320, true)
	if err != nil {
		return err
	}

	for _, frame := range frames {
		gocv.Resize(frame, &frame, image.Point{
			X: 480,
			Y: 320,
		}, 0, 0, gocv.InterpolationLinear)
		writer.Write(frame)
	}

	writer.Close()
	sendCh <- fname

	return nil
}

func (p *Processor) Process(ctx context.Context, path string, output chan string) error {
	capt, err := gocv.VideoCaptureFile(path)
	defer capt.Close()

	if err != nil {
		return err
	}

	log.Println("cap", capt)

	if err != nil {
		log.Panicln("panic", err)
	}

	img := gocv.NewMat()
	defer img.Close()

	imgSubRect := gocv.NewMat()
	defer imgSubRect.Close()

	imgDelta := gocv.NewMat()
	defer imgDelta.Close()

	imgThresh := gocv.NewMat()
	defer imgThresh.Close()

	mog2 := gocv.NewBackgroundSubtractorMOG2()
	defer mog2.Close()

	if ok := capt.Read(&img); !ok {
		log.Printf("Cannot read src\n")
		return ErrCantReadImg
	}

	fmt.Printf("Start reading device: \n")

	counter := 0

	movingSeq := make([]gocv.Mat, 0)

	for {
		counter++
		log.Printf("process frame %d of %s", counter, path)
		if ok := capt.Read(&img); !ok {
			fmt.Printf("Device closed: \n")
			return nil
		}
		if img.Empty() {
			continue
		}

		region := img.Region(image.Rect(0, 0, 600, 550))
		region.CopyTo(&imgSubRect)

		gocv.GaussianBlur(imgSubRect, &imgSubRect, image.Pt(85, 85), 5.5, 5.5, gocv.BorderConstant)

		// first phase of cleaning up image, obtain foreground only
		mog2.Apply(imgSubRect, &imgDelta)

		// remaining cleanup of the image to use for finding contours.
		// first use threshold
		gocv.Threshold(imgDelta, &imgThresh, 25, 255, gocv.ThresholdOtsu)

		// then dilate
		kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
		gocv.Dilate(imgThresh, &imgThresh, kernel)
		kernel.Close()

		// now find contours
		contours := gocv.FindContours(imgThresh, gocv.RetrievalList, gocv.ChainApproxSimple)

		csize := contours.Size()
		var numObj = 0
		for i := 0; i < csize; i++ {
			area := gocv.ContourArea(contours.At(i))
			if area < p.config.MinimumArea {
				continue
			}
			log.Println("contours at i ", contours.At(i).Size())
			numObj++

			rect := gocv.BoundingRect(contours.At(i))
			gocv.Rectangle(&region, rect, color.RGBA{0, 0, 255, 0}, 2)
			gocv.Rectangle(&region, image.Rect(0, 0, region.Rows(), region.Cols()), color.RGBA{255, 255, 255, 1}, 1)
		}

		contours.Close()

		if numObj > 0 {
			movingSeq = append(movingSeq, img.Clone())
		} else {
			if len(movingSeq) > 15 {
				go func(seq []gocv.Mat) {
					subrecord(ctx, seq, output)
				}(movingSeq)
			}
			movingSeq = make([]gocv.Mat, 0)
		}
	}

	return nil
}
