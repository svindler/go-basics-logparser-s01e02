package logparser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var httpLogLineRe = regexp.MustCompile(`^.* - (\d+.\d+.\d+T\d+:\d+:\d+\.\d+)`)

func ParseHttpLine(line string) (time.Time, error) {

	result := httpLogLineRe.FindSubmatch([]byte(line))

	// if we have errors
	if result == nil {
		return time.Time{}, fmt.Errorf("cannot find the timestamp in the line '%v'", line)
	}

	// parse the time
	parsedTime, err := time.Parse("2006-01-02T15:04:05", string(result[1]))
	if err != nil {
		return time.Time{}, fmt.Errorf("while parsing time: %v", err)
	}

	return parsedTime, nil
}

func CheckIfTimeInBetween(start, end, eventTime time.Time) bool {
	return eventTime.After(start) && eventTime.Before(end)
}




func CheckHttpFileForLines(filename string) ([]LogLine, error) {
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
		lineParsedTime, err := ParseHttpLine(lineText)
		if err != nil {
			fmt.Println("ERROR:", err)
			continue
		}

		logLines = append(logLines, LogLine{
			TimeStamp: lineParsedTime,
			Text:      lineText,
			Filename:  filename,
		})


	}

	return logLines, nil
}


type HttpAccessLineParser struct {
	FileGlob string
}



func (h *HttpAccessLineParser) Process(start, end time.Time) ([]LogLine, error) {

	// find relevant files
	matches, err := filepath.Glob(h.FileGlob)
	if err != nil {
		return nil, fmt.Errorf("while trying to glob '%v': %v", h.FileGlob, err)
	}



	logLinesMatched := make([]LogLine, 0)

	for _, match := range matches {
		fmt.Println("file: ", match)

		// find all lines
		logLines, err := CheckHttpFileForLines(match)
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
