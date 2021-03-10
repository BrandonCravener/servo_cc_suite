package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/tarm/serial"
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

	validRates := map[int]bool{300: true, 600: true, 1200: true, 2400: true, 4800: true, 9600: true, 14400: true, 19200: true, 28800: true, 31250: true, 38400: true, 57600: true, 115200: true}

	serialPort := flag.String("port", "COM1", "COM port the Arduino is attached to.")
	songFile := flag.String("file", "Default.song", "Path to the song file.")
	baudRate := flag.Int("baudRate", 115200, "Baud rate that the Arduino reads at.")
	doneWait := flag.Bool("doneWait", true, "Whether to wait for the DONE signal from the Arduino.")
	readBack := flag.Bool("readBack", false, "Read back the serial data being sent back from the Arduino.")

	flag.Parse()

	if !strings.HasPrefix(*serialPort, "COM") {
		fmt.Println("Invalid Windows Serial Port it should start with COM")
		os.Exit(1)
	}
	if !strings.HasSuffix(*songFile, ".song") {
		fmt.Println("Invalid song file all song files should end with '.song'")
		os.Exit(1)
	}
	if !validRates[*baudRate] {
		fmt.Println("Invalid baud rate, https://www.arduino.cc/en/Reference/SoftwareSerialBegin")
		os.Exit(1)
	}

	fmt.Printf("%v: Flags Read\n", time.Since(start))

	serialConfig := &serial.Config{Name: *serialPort, Baud: *baudRate, Size: 8}
	serial, err := serial.OpenPort(serialConfig)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer serial.Close()
	fmt.Printf("%v: Serial Connected\n", time.Since(start))

	file, err := os.Open(*songFile)
	if err != nil {
		panic(err)
	}

	if *doneWait {
		for {
			fmt.Printf("%v: Waiting Reset Confirmation...\n", time.Since(start))
			time.Sleep(5 * time.Second)
			buf := make([]byte, 128)
			n, err := serial.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			result := string(buf[:n])
			if strings.ContainsAny(result, "DONE") {
				break
			}
		}
	}
	fmt.Printf("%v: Arduino Reset Complete\n", time.Since(start))

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

		if *readBack {
			n, err := serial.Write([]byte(strconv.Itoa(servoNumber)))
			if err != nil {
				panic(err)
			}
			buf := make([]byte, 128)
			n, err = serial.Read(buf)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%q", n)
		} else {
			serial.Write([]byte(strconv.Itoa(servoNumber)))
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
