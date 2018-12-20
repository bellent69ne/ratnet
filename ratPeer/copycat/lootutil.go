// package lootutil contains utility functions
// for web loot operations
package copycat

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	//"os"
	"strconv"
	"time"
)

const (
	KB    = 1024
	MB    = KB * KB
	GB    = MB * MB
	CHUNK = 16 * KB
)

/***************YOU SHOULD RECONSIDER THIS CODE******************/
// inspect determines support for partial requests and size of file
func Inspect(url *string) (size int, err error) {
	// make HEAD request
	resp, err := http.Head(*url)
	// Handle errors if any
	if err != nil {
		return 0, errors.Wrap(err, "Failed making http head request")
	}
	// Check support for partial requests
	_, ok := resp.Header["Accept-Ranges"]
	if !ok {
		err = fmt.Errorf("%s doesn't support partial requests...", *url)
		return 0, err
	}
	fileSize, _ := resp.Header["Content-Length"]
	// Convert file size from string to int
	size, err = strconv.Atoi(fileSize[0])
	// Probably need to wrap err here////////////////////
	return
}

// Create http request
func createRequest(method, block string, url *string) (
	req *http.Request, err error) {
	req, err = http.NewRequest("GET", *url, nil)
	// if error happened, return immediately
	if err != nil {
		return nil, errors.Wrap(err, "Failed creating http request")
	}
	// Add "Range" header to http request
	req.Header.Add("Range", block)

	return
}

type DataChunk struct {
	data []byte
	time.Duration
	err error
}

// Try returning response body
// lootChunk downloads piece of file
func LootChunk(url *string, startAt int64, stream chan []byte) {
	block := fmt.Sprintf("bytes=%v-%v", startAt, startAt+(CHUNK-1))

	req, err := createRequest("GET", block, url)
	// Handle errors if any
	if err != nil {
		//stream <- DataChunk{nil, 0, err}
		stream <- nil
		return
	}

	// Create http client
	client := &http.Client{}
	// start measuring time
	//start := time.Now()
	// Make request
	resp, err := client.Do(req)
	// if error happened, send nil data and error to the channel
	if err != nil {
		//stream <- DataChunk{nil, 0, err}
		stream <- nil
		return
	}
	defer resp.Body.Close()
	// Read all bytes from the body of response
	body, err := ioutil.ReadAll(resp.Body)
	//elapsed := time.Now().Sub(start)
	// Send that body to a channel
	stream <- body //DataChunk{body, elapsed, err}
}

// Loot downloads pieces of file from url
// and assembles them into one file simultaneously
/*func Loot(url *string) error {
	size, err := Inspect(url)

	if err != nil {
		return err
	}

	stream := make(chan DataChunk)

	var nextChunk int64

	fi, err := os.Stat(Filename(url))
	switch {
	case err == nil:
		nextChunk = fi.Size()
	case os.IsNotExist(err):
		nextChunk = 0
	case err != nil && !os.IsNotExist(err):
		return err
	}

	fmt.Printf("\nFile: %s\n", Filename(url))
	for nextChunk != int64(size) {
		go LootChunk(url, nextChunk, stream)
		err := WriteToFile(Filename(url), stream)

		if err != nil {
			return err
		}

		nextChunk += CHUNK
		// ALL THIS CODE STILL REQUIERS MODIFICATION
		if nextChunk > int64(size) {
			nextChunk = int64(size)
		}
		//fmt.Printf("\telapsed %v\t%v\n", elapsed, int(elapsed))
		printStatus(nextChunk, size, elapsed)
	}
	fmt.Printf("\n\n")

	return err
}*/
