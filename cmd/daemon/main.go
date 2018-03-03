package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/nikita-vanyasin/go-web-course/common"
	"github.com/nikita-vanyasin/go-web-course/ffmpeg"
	"github.com/nikita-vanyasin/go-web-course/handlers"
	"github.com/nikita-vanyasin/go-web-course/video"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func taskProvider(stopChan chan struct{}, context *handlers.IsoContext) <-chan *video.Item {
	tasksChan := make(chan *video.Item)
	go func() {
		for {
			select {
			case <-stopChan:
				close(tasksChan)
				return
			default:
			}

			item, err := context.VideoRepository.GetUnprocessedItem()
			if err != nil {
				log.Error(err)
				close(tasksChan)
				return
			}

			if item != nil {
				log.Printf("got the task %v\n", item)

				item.Status = video.StatusProcessing
				err = context.VideoRepository.Update(item)
				if err != nil {
					log.Error(err)
					close(tasksChan)
					return
				}

				tasksChan <- item
			} else {
				time.Sleep(2 * time.Second)
			}
		}
	}()
	return tasksChan
}

func runTaskProvider(stopChan chan struct{}, context *handlers.IsoContext) <-chan *video.Item {
	resultChan := make(chan *video.Item)
	stopTaskProviderChan := make(chan struct{})
	taskProviderChan := taskProvider(stopTaskProviderChan, context)
	onStop := func() {
		stopTaskProviderChan <- struct{}{}
		close(resultChan)
	}
	go func() {
		for {
			select {
			case <-stopChan:
				onStop()
				return
			case task := <-taskProviderChan:
				select {
				case <-stopChan:
					onStop()
					return
				case resultChan <- task:
				}
			}
		}
	}()
	return resultChan
}

func worker(tasksChan <-chan *video.Item, name int, context *handlers.IsoContext) {

	handleError := func(err error, item *video.Item) bool {
		if err != nil {
			log.Printf("while processing %v occurred: %v", item, err)
			item.Status = video.StatusError
			anotherErr := context.VideoRepository.Update(item)
			if anotherErr != nil {
				log.Error(anotherErr) // TODO: stop all workers???
			}
		}
		return err != nil
	}

	log.Printf("start worker %v\n", name)
	for item := range tasksChan {
		log.Printf("start handle item %v on worker %v\n", item, name)

		duration, err := ffmpeg.GetVideoDuration(context.VideoStorage.GetFilePath(item))
		if handleError(err, item) {
			return
		}

		err = ffmpeg.CreateVideoThumbnail(context.VideoStorage.GetFilePath(item), context.VideoStorage.GetThumbnailPath(item), int64(duration)/2)
		if handleError(err, item) {
			return
		}

		item.Duration = int64(duration)
		item.Status = video.StatusReady
		err = context.VideoRepository.Update(item)
		if err != nil {
			log.Error(err) // TODO: stop all workers???
		}

		log.Printf("end handle item %v on worker %v\n", item, name)
	}
	log.Printf("stop worker %v\n", name)
}

func runWorkerPool(stopChan chan struct{}, context *handlers.IsoContext) *sync.WaitGroup {
	var wg sync.WaitGroup
	tasksChan := runTaskProvider(stopChan, context)
	for i := 0; i < 3; i++ {
		go func(i int) {
			wg.Add(1)
			worker(tasksChan, i, context)
			wg.Done()
		}(i)
	}
	return &wg
}

func main() {
	file := common.SetupLogging("daemon")
	defer file.Close()

	envSettings := common.GetEnvSettings()
	isoContext := handlers.CreateContext(envSettings)
	defer isoContext.Shutdown()

	stopChan := make(chan struct{})

	wg := runWorkerPool(stopChan, isoContext)

	killChan := getKillSignalChan()
	waitForKillSignal(killChan)
	stopChan <- struct{}{}
	wg.Wait()
}

func getKillSignalChan() chan os.Signal {
	osKillSignalChan := make(chan os.Signal, 1)
	signal.Notify(osKillSignalChan, os.Kill, os.Interrupt, syscall.SIGTERM)
	return osKillSignalChan
}

func waitForKillSignal(killSignalChan chan os.Signal) {
	killSignal := <-killSignalChan
	switch killSignal {
	case os.Interrupt:
		log.Info("got SIGINT...")
	case syscall.SIGTERM:
		log.Info("got SIGTERM...")
	}
}
