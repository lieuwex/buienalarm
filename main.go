package main

import (
	"buienalarm/location"
	"fmt"
	"os"
)

func printUsage() {
	fmt.Fprintf(os.Stderr, "usage: %s [-24] <location>\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
	}

	var locstr string
	fullDay := false
	if len(os.Args) == 2 {
		locstr = os.Args[1]
	} else if len(os.Args) == 3 {
		switch os.Args[1] {
		case "-2":
			fullDay = false
		case "-24":
			fullDay = true

		default:
			printUsage()
		}

		locstr = os.Args[2]
	}

	loc, err := location.Geocode(locstr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error during geocode, maybe unknown location?")
		os.Exit(1)
	}

	info, err := func() ([]timepoint, error) {
		if fullDay {
			return fetchFullDay(loc.Lat, loc.Lng)
		}

		return fetchTwoHours(loc.Lat, loc.Lng)
	}()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't retrieve forecast data from Buienradar: %s\n", err)
		os.Exit(1)
	}

	for _, point := range info {
		fmt.Printf(
			"%02d:%02d: %.1fmm/u\n",
			point.Time.Hour(),
			point.Time.Minute(),
			point.Precipation,
		)
	}
}
