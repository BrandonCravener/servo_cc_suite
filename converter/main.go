package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func convertToBin(midiFile, miditonesBinFile string) {
	midiBinCmd := exec.Command("miditones.exe", "-b", "-s2", "-delaymin=10", "-notemin=100", midiFile)
	err := midiBinCmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if _, err := os.Stat(miditonesBinFile); os.IsNotExist(err) {
		fmt.Println("Unable to find BIN file")
		os.Exit(1)
	}

}

func convertToTxt(midiFile, miditonesBinFile, miditonesTxtFile string) {
	midiTxtCmd := exec.Command("miditones_scroll.exe", strings.Split(miditonesBinFile, ".")[0])
	err := midiTxtCmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if _, err := os.Stat(miditonesTxtFile); os.IsNotExist(err) {
		fmt.Println("Unable to find TXT file")
		os.Exit(1)
	}
	if err := os.Remove(miditonesBinFile); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	progressBar := pb.Simple.Start(4)
	guitarStrings := map[string]int{
		"4E": 0,
		"B":  1,
		"G":  2,
		"D":  3,
		"A":  4,
		"3E": 5,
	}

	inputMidi := flag.String("file", "DEFAULT", "Path to MIDI file.")
	zeroDelay := flag.Int("zeroDelay", 1, "Delay between notes that get played at the same time.")

	flag.Parse()

	if !strings.HasSuffix(*inputMidi, ".mid") {
		fmt.Println("Invalid MIDI file ensure it has a .mid extension")
		os.Exit(1)
	}
	if _, err := os.Stat("miditones.exe"); os.IsNotExist(err) {
		fmt.Println("Unable to find miditones")
		os.Exit(1)
	}

	progressBar.Increment()

	miditonesBinFile := strings.ReplaceAll(filepath.Base(*inputMidi), ".mid", ".bin")
	miditonesTxtFile := strings.ReplaceAll(filepath.Base(*inputMidi), ".mid", ".txt")
	songFilePath := strings.ReplaceAll(filepath.Base(*inputMidi), ".mid", ".song")

	convertToBin(*inputMidi, miditonesBinFile)

	progressBar.Increment()

	convertToTxt(*inputMidi, miditonesBinFile, miditonesTxtFile)

	progressBar.Increment()

	txtFile, err := os.Open(miditonesTxtFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	songFile, _ := os.Create(songFilePath)

	scanner := bufio.NewScanner(txtFile)
	tableRead := false
	noteRegex, _ := regexp.Compile(`(?:(3|4)E)|([A-z])`)

	for scanner.Scan() {
		text := scanner.Text()
		if len(text) == 0 {
			continue
		}
		if strings.HasSuffix(text, " used.") {
			tableRead = false
		}
		// TODO: Add dedupe output to .song format
		if tableRead {
			// Splice data to notes and delay
			splitData := strings.Fields(text)
			for i := 0; i < len(splitData); i++ {
				if strings.HasSuffix(splitData[i], ":") {
					splitData = splitData[1:i]
					break
				}
			}

			// Add notes to de duplicated data
			var dedupeData []int
			for i := 0; i < len(splitData)-1; i++ {
				servoNum := guitarStrings[noteRegex.FindString(splitData[i])]
				if !intInSlice(servoNum, dedupeData) {
					dedupeData = append(dedupeData, servoNum)
				}
			}
			// Extract delay and convert to int ms
			delayMS, _ := strconv.ParseFloat(splitData[len(splitData)-1], 8)
			delayMS = delayMS * 1000

			// fmt.Printf("%v|%v\n", text, delayMS)
			// fmt.Println(dedupeData)
			// Convert array into song file lines
			if len(dedupeData) > 0 {
				for i, note := range dedupeData {
					if i == len(dedupeData)-1 {
						songFile.WriteString(fmt.Sprintf("%v:%v\n", note, delayMS))
					} else {
						songFile.WriteString(fmt.Sprintf("%v:%v\n", note, *zeroDelay))
					}
				}
			} else {
				songFile.WriteString(fmt.Sprintf("99:%v\n", delayMS))
			}

		}
		if strings.HasSuffix(text, "bytestream code") {
			tableRead = true
		}
	}
	progressBar.Increment()

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Clean up
	if err := txtFile.Close(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := songFile.Close(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := os.Remove(miditonesTxtFile); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	progressBar.Finish()
}
