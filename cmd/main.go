package main

import (
	"blob/internal/app"
	"blob/internal/delivery/telegram"
	"blob/internal/input"
	"blob/internal/video_processor"
	"context"
	"log"
	"os"
	"sync"
	"time"
)

func init() {
	//format.RegisterAll()
}

func main() {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cfg, err := app.LoadConfig("./configs")

	if err != nil {
		log.Fatal(err)
	}

	logging := log.New(os.Stdout, "", log.LstdFlags)

	tg, err := telegram.NewTelegram(&cfg.Telegram)
	if err != nil {
		log.Panicln(err)
	}
	_ = tg.AddCommandHandler("help", "Список команд", telegram.HandleHelp)
	_ = tg.AddCommandHandler("start", "Начать уведомлять вас о признаках движения", telegram.HandleSub)
	_ = tg.AddCommandHandler("stop", "Перестать уведомлять", telegram.HandleUnsub)

	err = tg.HandleSubscribers()

	filesCh := make(chan string)
	defer close(filesCh)

	tg.SendToSubs(filesCh)

	if err != nil {
		log.Println("error handle subscribers")
	}

	procConf := video_processor.ProcessorConfig{MinimumArea: 2000}
	proc := video_processor.NewProcessor(procConf)

	wg := &sync.WaitGroup{}

	grabberConf := &input.GrabberConf{
		PlaylistURL: "http://109.106.138.159:3568/27329254-5b18-11eb-ae93-0242ac130002/cam3/f/index.m3u8",
		UpdateEvery: 30 * time.Second,
	}
	grabber := input.NewGrabber(grabberConf, logging)

	urls := grabber.GrabURLS(ctx, grabber.LoadPlaylistContent)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for job := range urls {
			job := job
			go func() {
				err := proc.Process(ctx, job, filesCh)
				if err != nil {
					log.Println(err)
				}
			}()
		}

	}()
	//go func() {
	//	if err := proc.Process(ctx, "./files/43950dcc-ec93-420d-a6c0-da0885039f59.ts", filesCh); err != nil {
	//		log.Fatal(err)
	//	}
	//}()

	wg.Wait()

}
