package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type timepoint struct {
	Time   time.Time
	Percip float64
}

func parseLine(line string) (timepoint, error) {
	fields := strings.Split(line, "|")

	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return timepoint{}, err
	}

	time, err := time.Parse("15:04", fields[1])
	if err != nil {
		return timepoint{}, err
	}

	return timepoint{
		Time:   time,
		Percip: math.Pow(10, (float64(val)-109)/32),
	}, nil
}

func fetch(latitude float64, longitude float64) ([]timepoint, error) {
	url := fmt.Sprintf("https://gpsgadget.buienradar.nl/data/raintext?lat=%f&lon=%f", latitude, longitude)
	resp, err := http.Get(url)
	if err != nil {
		return []timepoint{}, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []timepoint{}, err
	}
	lines := strings.Split(string(bytes), "\n")

	var res []timepoint
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		point, err := parseLine(line)
		if err != nil {
			return []timepoint{}, err
		}
		res = append(res, point)
	}
	return res, nil
}
