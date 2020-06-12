package main

import (
	"flag"
	"fmt"
	"logparser"
	"os"
	"time"
)

type jsonLogLine struct {
	Ts string `json:"ts"`
}

func main() {

	startStringRef := flag.String("start", "2020-06-11T13:48:00.000", "The start time to look for in logs")
	endStringRef := flag.String("end", "2020-06-19T13:48:00.000", "The end time to look for in logs")

	flag.Parse()


	if startStringRef == nil || endStringRef == nil {
		fmt.Println("Missing start and / or end time")
		os.Exit(-1)
	}



	parser := logparser.CombinedParser{
		Parsers: []logparser.LogFileParser{
			&logparser.HttpAccessLineParser{"/Users/gyulalaszlo/Documents/TableauLogs/logs/httpd/access.*.log"},
			&logparser.JsonLogLineParser{"/Users/gyulalaszlo/Documents/TableauLogs/logs/vizqlserver/nativeapi_vizqlserver*.txt"},
		},
	}

	// Check times validity
	//start, err := time.Parse("2006-01-02T15:04:05", "2020-06-11T13:48:00.000")
	start, err := time.Parse("2006-01-02T15:04:05", *startStringRef)
	if err != nil {
		fmt.Println("Malformed start time:", err)
		os.Exit(-1)
	}

	end, err := time.Parse("2006-01-02T15:04:05", *endStringRef)
	if err != nil {
		fmt.Println("Malformed end time:", err)
		os.Exit(-1)
	}

	logLines, err := parser.Process(start, end)

	if err != nil {
		panic(err)
	}

	for _, line := range logLines {

		if logparser.CheckIfTimeInBetween(start, end, line.TimeStamp) {
			fmt.Println(line)
		}
	}

}
