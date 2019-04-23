package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/aerogo/http/client"
	"github.com/akyoto/color"
)

// Crawler ...
type Crawler struct {
	headers client.Headers
	ticker  *time.Ticker
	tasks   chan *Task
	wg      *sync.WaitGroup
}

// Task ...
type Task struct {
	URL         string
	Destination string
	Raw         bool
}

// Queue queues up a task.
func (crawler *Crawler) Queue(task *Task) {
	crawler.wg.Add(1)
	crawler.tasks <- task
}

// Wait waits until all tasks have been completed.
func (crawler *Crawler) Wait() {
	crawler.wg.Wait()
}

// Download page contents to file.
func (crawler *Crawler) work(task *Task) {
	response, err := client.Get(task.URL).Headers(crawler.headers).End()

	if err != nil {
		fmt.Println(color.RedString(task.URL), err)
		return
	}

	if response.StatusCode() != http.StatusOK {
		fmt.Println(color.RedString(task.URL), response.StatusCode())
		return
	}

	var data []byte

	if task.Raw {
		data = response.RawBytes()
	} else {
		data = response.Bytes()
	}

	fmt.Println(color.GreenString(task.URL), len(data), "bytes")
	ioutil.WriteFile(task.Destination, data, 0644)
}

// New ...
func New(headers client.Headers, delayBetweenRequests time.Duration, tasksBufferLength int) *Crawler {
	crawl := &Crawler{
		headers: headers,
		ticker:  time.NewTicker(delayBetweenRequests),
		tasks:   make(chan *Task, tasksBufferLength),
		wg:      &sync.WaitGroup{},
	}

	go func() {
		for task := range crawl.tasks {
			crawl.work(task)
			crawl.wg.Done()
			<-crawl.ticker.C
		}
	}()

	return crawl
}
