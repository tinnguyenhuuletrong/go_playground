package play_sync

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func Play_Error_Group_With_Capping() {
	g := new(errgroup.Group)
	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://www.somestupidname.com/",
		"http://haha-fake/",
	}

	g.SetLimit(2)

	results := make([]bool, len(urls))

	for i, url := range urls {
		// Launch a goroutine to fetch the URL.
		url := url // https://golang.org/doc/faq#closures_and_goroutines
		jobIndex := i
		g.Go(func() error {
			// Fetch the URL.
			log.Println("begin fetch:", url)

			resp, err := http.Get(url)
			if err == nil {
				resp.Body.Close()
			}
			defer log.Println("end fetch:", url, err)
			results[jobIndex] = err == nil
			return err
		})
	}

	// Wait for all HTTP fetches to complete.
	g.Wait()

	fmt.Println("Successfully fetched all URLs.", results)

}
