package crawler

import "time"

// Crawler ...
type Crawler struct {
	Ticker *time.Ticker
}

// Task ...
type Task struct {
	URL         string
	Destination string
}

// Queue queues a given URL to be downloaded to the given file destination.
func (crawler *Crawler) Queue(task *Task) {

}

// New ...
func New(delayBetweenRequests time.Duration) *Crawler {
	return &Crawler{
		Ticker: time.NewTicker(delayBetweenRequests),
	}
}
