package crawler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/aerogo/http/client"
	"github.com/akyoto/color"
)

// Crawler is a web crawler that accepts URLs and tries to fetch them.
type Crawler struct {
	headers client.Headers
	ticker  *time.Ticker
	tasks   chan *Task
	wg      *sync.WaitGroup
}

// Task represents a single URL fetch task.
type Task struct {
	URL         string
	Destination string
	Raw         bool
}

// Queue queues up a task.
func (crawler *Crawler) Queue(task *Task) error {
	if task.URL == "" {
		return errors.New("Task is missing a URL")
	}

	if task.Destination == "" {
		return fmt.Errorf("Task '%s' is missing a destination", task.URL)
	}

	crawler.wg.Add(1)
	crawler.tasks <- task
	return nil
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
		data = response.Raw()
	} else {
		data = response.Bytes()
	}

	fmt.Println(color.GreenString(task.URL), len(data), "bytes")
	err = ioutil.WriteFile(task.Destination, data, 0644)

	if err != nil {
		fmt.Println(color.RedString(task.URL), err)
		return
	}
}

// New creates a new crawler.
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
