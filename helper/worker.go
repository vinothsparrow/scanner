package helper

import (
	log "github.com/sirupsen/logrus"
	"github.com/vinothsparrow/scanner/model"
)

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
				log.Printf("worker%d: Received work request, delaying for %f seconds\n", w.Id, work.Id)

				log.Printf("worker%d: Hello, %s!\n", w.Id, work.Id)

			case <-w.QuitChan:
				// We have been asked to stop.
				log.Printf("worker%d stopping\n", w.Id)
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
