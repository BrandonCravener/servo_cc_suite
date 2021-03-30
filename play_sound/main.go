package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

func lineCounter(r io.Reader) (int, error) {
	var count int
	const lineBreak = '\n'

	buf := make([]byte, bufio.MaxScanTokenSize)

	for {
		bufferSize, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}

		var buffPosition int
		for {
			i := bytes.IndexByte(buf[buffPosition:], lineBreak)
			if i == -1 || bufferSize == buffPosition {
				break
			}
			buffPosition += i + 1
			count++
		}
		if err == io.EOF {
			break
		}
	}

	return count, nil
}

func main() {
	start := time.Now()

	songFile := flag.String("file", "Default.song", "Path to the song file.")

	flag.Parse()
	if !strings.HasSuffix(*songFile, ".song") {
		fmt.Println("Invalid song file all song files should end with '.song'")
		os.Exit(1)
	}

	fmt.Printf("%v: Flags Read\n", time.Since(start))

	file, err := os.Open(*songFile)
	if err != nil {
		panic(err)
	}

	var files []*os.File
	var buffers []*beep.Buffer
	var formats []beep.Format

	for i := 0; i <= 5; i++ {
		temp, err := os.Open(fmt.Sprintf("sounds/%v.wav", i))
		if err != nil {
			panic(err)
		}
		files = append(files, temp)
	}
	for _, f := range files {
		streamer, format, err := wav.Decode(f)
		if err != nil {
			panic(err)
		}
		buffer := beep.NewBuffer(format)
		buffer.Append(streamer)
		streamer.Close()
		buffers = append(buffers, buffer)
		formats = append(formats, format)
	}

	speaker.Init(formats[0].SampleRate, formats[0].SampleRate.N(time.Millisecond*50))

	lines, _ := lineCounter(file)
	progressBar := pb.StartNew(lines)
	file.Close()

	file, err = os.Open(*songFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// previousCommandTime := time.Now()
	for scanner.Scan() {
		progressBar.Increment()
		command := scanner.Text()
		split := strings.Split(command, ":")

		servoNumber, err := strconv.Atoi(split[0])
		if err != nil {
			continue
		}
		delay, err := strconv.Atoi(split[1])
		if err != nil {
			continue
		}

		switch servoNumber {
		case 0:
			temp := buffers[0].Streamer(0, buffers[0].Len())
			speaker.Play(temp)
		case 1:
			temp := buffers[1].Streamer(0, buffers[1].Len())
			speaker.Play(temp)
		case 2:
			temp := buffers[2].Streamer(0, buffers[2].Len())
			speaker.Play(temp)
		case 3:
			temp := buffers[3].Streamer(0, buffers[3].Len())
			speaker.Play(temp)
		case 4:
			temp := buffers[4].Streamer(0, buffers[4].Len())
			speaker.Play(temp)
		case 5:
			temp := buffers[5].Streamer(0, buffers[5].Len())
			speaker.Play(temp)
		}

		// fmt.Printf("%s|%s: %v, %vms\n", time.Since(start), time.Since(previousCommandTime), servoNumber, delay)
		// previousCommandTime = time.Now()

		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	progressBar.Finish()
}
