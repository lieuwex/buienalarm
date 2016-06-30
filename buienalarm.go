package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

type weatherInfo struct {
	fullDay bool
	Start   float64   `json:"start"`
	Precip  []float64 `json:"precip"`
	Raw     []float64 `json:"raw"`
}

func (info weatherInfo) getAmount() [25]float64 {
	var result [25]float64

	if info.fullDay {
		for i, val := range info.Raw {
			result[i] = val
		}
	} else {
		for i, val := range info.Precip {
			result[i] = math.Pow(10, float64((val-109)/32))
		}
	}

	return result
}

func parseString(s string, fullDay bool) weatherInfo {
	x := strings.Split(s, " ")
	raw := strings.Replace(x[len(x)-1], ";", "", -1)

	var data weatherInfo
	err := json.Unmarshal([]byte(raw), &data)
	if err != nil {
		fmt.Println("Couldn't parse JSON. ")
		panic(err)
	}

	data.fullDay = fullDay
	return data
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "usage: %s [-24] <location>\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
	}

	var location string
	fullDay := false
	if len(os.Args) == 2 {
		location = os.Args[1]
	} else if len(os.Args) == 3 {
		switch os.Args[1] {
		case "-2":
			fullDay = false
		case "-24":
			fullDay = true

		default:
			printUsage()
		}

		location = os.Args[2]
	}

	resp, err := http.Get("http://www.buienalarm.nl/location/" + location)
	if err != nil {
		fmt.Println("Couldn't retreive forcast data from Buienalarm. ")
		panic(err)
	}
	defer resp.Body.Close()

	matchString := "locationdata['forecast']"
	timeInterval := 5 * time.Minute
	if fullDay {
		matchString = "var precip_daily"
		timeInterval = time.Hour
	}

	body, _ := ioutil.ReadAll(resp.Body)
	for _, line := range strings.Split(string(body), "\n") {
		if strings.Contains(line, matchString) {
			data := parseString(line, fullDay)

			t := time.Unix(int64(data.Start), 0)
			for _, val := range data.getAmount() {
				fmt.Printf("%02d:%02d: %.1fmm/u\n", t.Hour(), t.Minute(), val)
				t = t.Add(timeInterval)
			}

			return
		}
	}

	fmt.Fprintf(os.Stderr, "location not found: '%s'\n", location)
	os.Exit(1)
}
