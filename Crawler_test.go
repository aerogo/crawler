package crawler_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aerogo/crawler"
	"github.com/aerogo/http/client"
)

func Test(t *testing.T) {
	defer os.Remove("test-0.html")
	defer os.Remove("test-1.html")

	n := 2
	c := crawler.New(client.Headers{}, 100*time.Millisecond, n)

	for i := 0; i < n; i++ {
		_ = c.Queue(&crawler.Task{
			URL:         "https://eduardurbach.com",
			Destination: fmt.Sprintf("test-%d.html", i),
			Raw:         i%2 == 0,
		})
	}

	c.Wait()
}
