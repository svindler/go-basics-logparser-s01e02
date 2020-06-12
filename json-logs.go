package logparser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type JsonLogLineParser struct {
	FileGlob string
}


func (h *JsonLogLineParser) Process(start, end time.Time) ([]LogLine, error) {

	// find relevant files
	matches, err := filepath.Glob(h.FileGlob)
	if err != nil {
		return nil, fmt.Errorf("while trying to glob '%v': %v", h.FileGlob, err)
	}

	logLinesMatched := make([]LogLine, 0)

	for _, match := range matches {
		fmt.Println("file: ", match)

		// find all lines
		logLines, err := h.parseFile(match)
		if err != nil {
			return nil, fmt.Errorf("while trying to parse log lines from file '%v': %v", match, err)
		}


		// do the checking
		for _, line := range logLines {
			if CheckIfTimeInBetween(start, end, line.TimeStamp) {
				logLinesMatched = append(logLinesMatched, line)
			}
		}
	}

	return logLinesMatched, nil
}

type jsonLogLine struct {
	Ts string `json:"ts"`
}



func (h *JsonLogLineParser)parseFile(filename string) ([]LogLine, error) {
	// open the file
	file, err := os.Open(filename) // For read access.
	if err != nil {
		return nil, fmt.Errorf("while opening '%v' for reading: %v", filename, err)
	}
	// close the file
	defer file.Close()

	// create the empty container
	logLines := make([]LogLine, 0)

	// read line-by-line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineText := scanner.Text()

		var jsonLine jsonLogLine

		// parse json
		if err := json.Unmarshal([]byte(lineText), &jsonLine ); err != nil {
			continue
		}

		//
		parsedTime, err := time.Parse("2006-01-02T15:04:05", jsonLine.Ts)
		if err != nil {
			fmt.Println("ERROR:", err)
			continue
		}

		logLines = append(logLines, LogLine{
			TimeStamp: parsedTime,
			Text:      lineText,
			Filename:  filename,
		})


	}

	return logLines, nil
}
