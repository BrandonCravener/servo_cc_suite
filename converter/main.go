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
)

func convertToBin(midiFile, miditonesBinFile string) {
	midiBinCmd := exec.Command("miditones.exe", "-b", midiFile)
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

}

func main() {
	guitarStrings := map[string]int{
		"4E": 0,
		"B":  1,
		"G":  2,
		"D":  3,
		"A":  4,
		"3E": 5,
	}

	inputMidi := flag.String("file", "DEFAULT", "song_name.mid")

	flag.Parse()

	if !strings.HasSuffix(*inputMidi, ".mid") {
		fmt.Println("Invalid MIDI file ensure it has a .mid extension")
		os.Exit(1)
	}
	if _, err := os.Stat("miditones.exe"); os.IsNotExist(err) {
		fmt.Println("Unable to find miditones")
		os.Exit(1)
	}

	miditonesBinFile := strings.ReplaceAll(filepath.Base(*inputMidi), ".mid", ".bin")
	miditonesTxtFile := strings.ReplaceAll(filepath.Base(*inputMidi), ".mid", ".txt")

	fmt.Println("Converting to BIN")

	convertToBin(*inputMidi, miditonesBinFile)

	fmt.Println("Done.")
	fmt.Println("Converting to TXT")

	convertToTxt(*inputMidi, miditonesBinFile, miditonesTxtFile)

	fmt.Println("Done.")

	txtFile, err := os.Open(miditonesTxtFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer txtFile.Close()

	scanner := bufio.NewScanner(txtFile)
	tableRead := false
	noteRegex, _ := regexp.Compile(`(?:(3|4)E)|([A-z])`)

	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasSuffix(text, " used.") {
			tableRead = false
		}
		if len(text) == 0 {
			continue
		}
		// TODO: Add dedupe output to .song format
		if tableRead {
			splitData := strings.Fields(text)
			for i := 0; i < len(splitData); i++ {
				if strings.HasSuffix(splitData[i], ":") {
					splitData = splitData[1:i]
					break
				}
			}
			for i := 0; i < len(splitData)-1; i++ {
				fmt.Print(guitarStrings[noteRegex.FindString(splitData[i])])
				fmt.Print(" ")
			}
			delayMS, _ := strconv.ParseFloat(splitData[len(splitData)-1], 8)
			delayMS = delayMS * 1000
			fmt.Printf("%v ms\n", delayMS)
		}
		if strings.HasSuffix(text, "bytestream code") {
			tableRead = true
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
