package lib

import "log"

///////////////////////////////////////////////////////////////////////////////

// Represents a worker that can be added to the queue of workers
type Worker struct {
	ID          int
	Work        chan YouTubeVideo
	WorkerQueue chan chan YouTubeVideo
	QuitChan    chan bool
}

// NewWorker creates and returns a new Worker object. Its only argument is a
// channel that the worker can add itself to whenever it is done its work
func NewWorker(id int, workerQueue chan chan YouTubeVideo) Worker {
	// Create and return the worker.
	worker := Worker{
		ID:          id,
		Work:        make(chan YouTubeVideo),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
	}

	return worker
}

// This function "starts" the worker by starting a goroutine, that is an
// infinite "for-select" loop.
func (w *Worker) Start() {
	go func() {
		for {
			// Add w into the worker queue
			w.WorkerQueue <- w.Work

			select {
			case video := <-w.Work:
				// Receive a video to download
				log.Printf("worker%d: Started downloading video %s", w.ID, video.VideoId)
				var downloaded = make(chan bool)
				go video.Download(downloaded)
				value := <-downloaded
				if value == false {
					log.Printf("worker%d: Error in downloading video %s", w.ID, video.VideoId)
				}
				log.Printf("worker%d: Video %s downloaded", w.ID, video.VideoId)
			case <-w.QuitChan:
				log.Printf("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stops tells the worker to stop listening for work requests.
// Note that the worker will only stop *after* it has finished its work.
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

///////////////////////////////////////////////////////////////////////////////

var WorkerQueue chan chan YouTubeVideo

// StartDispatcher starts a number of n workers.
func StartDispatcher(nworkers int) {
	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue = make(chan chan YouTubeVideo, nworkers)

	// Now, create all of workers
	for i := 0; i < nworkers; i++ {
		log.Printf("Starting worker %v\n", i+1)
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case video := <-VideoQueue:
				worker := <-WorkerQueue
				worker <- video
			}
		}
	}()
}

///////////////////////////////////////////////////////////////////////////////
