package copycat

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func ttyWidth() (width int, err error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()

	if err != nil {
		return 0, errors.Wrap(err, "Failed calculating tty width")
	}

	out = out[:len(out)-1]
	strOut := string(out)
	splitted := strings.Split(strOut, " ")

	width, err = strconv.Atoi(splitted[1])
	return
}

func state(percent int) string {
	width, err := ttyWidth()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	totalLength := width * 35 / 100
	stateLength := totalLength * percent / 100

	state := make([]rune, totalLength+2)
	state[0] = '|'
	for i := 1; i <= totalLength; i++ {
		if i == stateLength {
			state[i] = '>'
		} else if i < stateLength {
			state[i] = '='
		} else {
			state[i] = ' '
		}
	}
	state[len(state)-1] = '|'
	return string(state)
}

func printStatus(nextChunk int64, size int64, elapsed time.Duration) {
	fileSize := int64(0)
	var strSize, strGot string
	switch {
	case size > KB && size < MB:
		{
			fileSize = size / KB
			strSize = fmt.Sprintf("%dkB", fileSize)
			strGot = fmt.Sprintf("%.2fkB", float64(nextChunk)/KB)
		}

	case size > MB && size < GB:
		{
			fileSize = size / MB
			strSize = fmt.Sprintf("%dmB", fileSize)
			strGot = fmt.Sprintf("%.2fmB", float64(nextChunk)/MB)
		}
	}
	speed := CHUNK / (elapsed / time.Millisecond)
	speed *= 1000
	speed /= 1024
	percent := int(float64(nextChunk) / float64(size) * 100)
	fmt.Printf("\r    %s %d%%    %s    %s    %dkB/s    ",
		state(percent), percent, strGot, strSize, int(speed))
}
