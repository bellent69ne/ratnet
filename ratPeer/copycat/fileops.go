package copycat

import (
	"github.com/pkg/errors"
	"os"
	"strings"
	//"time"
)

// filename is getting a filename
func Filename(url *string) string {
	result := strings.Split(*url, "/")

	return result[len(result)-1]
}

func accessFile(filename *string) (file *os.File, err error) {
	// Create file object to open an existing file or create new one
	if _, err = os.Stat(*filename); os.IsNotExist(err) {
		file, err = os.Create(*filename)
	} else {
		// Open an existing file and move to the end of the file
		file, err = os.OpenFile(*filename, os.O_WRONLY|os.O_APPEND, 0644)
	}
	return
}

// writeToFile writes all data from a channel to a file
func WriteToFile(filename string, stream chan []byte) error { //DataChunk) (

	// receive data from stream, handle any errors
	received := <-stream
	//if received.err != nil {
	//	return errors.Wrap(received.err, "Unable to get data from web")
	//}

	file, err := accessFile(&filename)
	if err != nil {
		return errors.Wrap(err, "Couldn't access the file")
	}
	defer file.Close()
	// Write all data to the file
	_, err = file.Write(received)
	if err != nil {
		return errors.Wrap(err, "Failed writing data to a file")
	}
	file.Sync()
	return nil //received.Duration, nil
}
