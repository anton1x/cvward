package main

import (
	"blob/internal"
	"blob/internal/delivery"
	"context"
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"log"
	"strconv"
)

func init() {
	//format.RegisterAll()
}

const MinimumArea = 2000

func main() {

	ctx := context.Background()

	tg, err := delivery.NewTelegram()
	if err != nil {
		log.Panicln(err)
	}

	err = tg.HandleSubscribers()

	filesCh := make(chan string)
	defer close(filesCh)

	tg.SendToSubs(filesCh)

	if err != nil {
		log.Println("error handle subscribers")
	}

	capt, err := gocv.VideoCaptureFile("./43950dcc-ec93-420d-a6c0-da0885039f59.ts")
	defer capt.Close()

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
		fmt.Printf("Cannot read src\n")
		return
	}

	vwriter, _ := gocv.VideoWriterFile("out.mp4", "H264", 25, img.Cols(), img.Rows(), true)
	defer vwriter.Close()

	status := "Ready"

	fmt.Printf("Start reading device: \n")

	counter := 0

	movingSeq := make([]gocv.Mat, 0)

	for {
		counter++
		log.Println("process frame ", counter)
		if ok := capt.Read(&img); !ok {
			fmt.Printf("Device closed: \n")
			return
		}
		if img.Empty() {
			continue
		}

		region := img.Region(image.Rect(400, 200, 600, 550))
		region.CopyTo(&imgSubRect)

		gocv.GaussianBlur(imgSubRect, &imgSubRect, image.Pt(85, 85), 5.5, 5.5, gocv.BorderConstant)

		status = "Ready"
		statusColor := color.RGBA{0, 255, 0, 0}

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
			if area < MinimumArea {
				continue
			}
			log.Println("contours at i ", contours.At(i).Size())
			numObj++
			//status = "Motion detected"
			//statusColor = color.RGBA{255, 0, 0, 0}
			//gocv.DrawContours(&img, contours, i, statusColor, 2)

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
					internal.Subrecord(ctx, seq, filesCh)
				}(movingSeq)
			}
			movingSeq = make([]gocv.Mat, 0)
		}

		gocv.PutText(&img, status, image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, statusColor, 2)
		gocv.PutText(&img, strconv.Itoa(numObj), image.Pt(img.Cols()-300, 20), gocv.FontItalic, 1.2, statusColor, 2)
		vwriter.Write(img)
		//vwriter2.Write(imgDelta)
		//vwriter3.Write(imgThresh)

		//outImg, err := img.ToImage()
		//path := "cv.png"
		//f, _ := os.Create(path)
		//png.Encode(f, outImg)

	}
}
