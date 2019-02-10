package helper

import (
	"errors"
	"time"

	cache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"github.com/vinothsparrow/scanner/model"
)

var Cache *cache.Cache

type Worker struct {
	Id          int
	Work        chan model.ScanRequest
	WorkerQueue chan chan model.ScanRequest
	QuitChan    chan bool
}

func NewWorker(id int, workerQueue chan chan model.ScanRequest) Worker {
	// Create, and return the worker.
	worker := Worker{
		Id:          id,
		Work:        make(chan model.ScanRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool)}

	return worker
}

func init() {
	Cache = cache.New(5*time.Hour, 24*time.Hour)
}

// This function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.
func (w *Worker) Start() {
	go func() {
		for {
			// Add ourselves into the worker queue.
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				// Receive a work request.
				log.Printf("worker%d: Received scan request, %s!\n", w.Id, work.Id)
				work.Status = "processing"
				Cache.Set(work.Id, &work, cache.DefaultExpiration)
				_, _ = GitScan(&work)
				work.Status = "done"
				Cache.Set(work.Id, &work, cache.DefaultExpiration)

			case <-w.QuitChan:
				// We have been asked to stop.
				log.Printf("worker%d stopping\n", w.Id)
				return
			}
		}
	}()
}

func GetScanRequest(id string) (*model.ScanRequest, error) {
	if req, found := Cache.Get(id); found {
		return req.(*model.ScanRequest), nil
	}
	return nil, errors.New("Scan not found")
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
