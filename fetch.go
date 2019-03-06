package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type timepoint struct {
	Time        time.Time
	Precipation float64
}

func fetchTwoHours(latitude float64, longitude float64) ([]timepoint, error) {
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

	parseLine := func(line string) (timepoint, error) {
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
			Time:        time,
			Precipation: math.Pow(10, (float64(val)-109)/32),
		}, nil
	}

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

func fetchFullDay(latitude float64, longitude float64) ([]timepoint, error) {
	url := fmt.Sprintf("https://graphdata.buienradar.nl/2.0/forecast/geo/rain24hour/?lat=%f&lon=%f", latitude, longitude)
	resp, err := http.Get(url)
	if err != nil {
		return []timepoint{}, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	var item struct {
		Forecasts []struct {
			Time        string  `json:"datetime"`
			Precipation float64 `json:"precipation"`
		} `json:"forecasts"`
	}
	if err := decoder.Decode(&item); err != nil {
		return []timepoint{}, err
	}

	var res []timepoint
	for _, item := range item.Forecasts {
		time, err := time.Parse("2006-01-02T15:04:05", item.Time)
		if err != nil {
			return []timepoint{}, err
		}

		res = append(res, timepoint{
			Time:        time,
			Precipation: item.Precipation,
		})
	}
	return res, nil
}
