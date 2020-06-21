package main

import (
	"flag"
	"fmt"
	"logparser"
	"os"
	"time"
)

/* type jsonLogLine struct {
	Ts string `json:"ts"`
}
*/

func main() {
	fmt.Println("Starting parser")

	startStringRef := flag.String("start", "2020-06-10T13:40:00.000", "The start time to look for in logs")
	endStringRef := flag.String("end", "2020-06-10T14:59:00.000", "The end time to look for in logs")

	flag.Parse()

	// Check here because I'm paranoid and these are user supplied pointers that will be dereferenced
	if startStringRef == nil || endStringRef == nil {
		fmt.Println("Missing start and / or end time")
		os.Exit(-1)
	}

	parser := logparser.CombinedParser{
		Parsers: []logparser.LogFileParser{
			&logparser.ControlLicenseServiceLogLineParser{"/Users/svr/Dev/TableauLogs/licenseservice/control_licenseservice_node1-0.log"},
			&logparser.ControlLicenseServiceLogLineParser{"/Users/svr/Dev/TableauLogs/licenseservice/control_licenseservice_node1-0.log.*"},
			&logparser.HttpdErrorLineParser{"/Users/svr/Dev/TableauLogs/httpd/error.log"},
			&logparser.HttpAccessLineParser{"/Users/svr/Dev/TableauLogs/httpd/access.*.log"},
			&logparser.JsonLogLineParser{"/Users/svr/Dev/TableauLogs/vizqlserver/nativeapi_vizqlserver*.txt"},
		},
	}

	// Check times validity
	// start, err := time.Parse("2006-01-02T15:04:05", "2020-06-11T13:48:00.000")
	start, err := time.Parse("2006-01-02T15:04:05", *startStringRef)
	if err != nil {
		fmt.Println("Malformed start time:", err)
		os.Exit(-1)
	}

	end, err := time.Parse("2006-01-02T15:05:05", *endStringRef)
	if err != nil {
		fmt.Println("Malformed end time:", err)
		os.Exit(-1)
	}

	logLines, err := parser.Process(start, end)

	if err != nil {
		panic(err)
	}

	for _, line := range logLines {
		fmt.Println(line)
	}

}
