package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"io"
)

type timepoint struct {
	Time        time.Time
	Precipation float64
}

func parseResponse(r io.ReadCloser) ([]timepoint, error){
	defer r.Close()

	decoder := json.NewDecoder(r)

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

func fetchTwoHours(latitude float64, longitude float64) ([]timepoint, error) {
	url := fmt.Sprintf("https://graphdata.buienradar.nl/2.0/forecast/geo/rain/?lat=%f&lon=%f", latitude, longitude)
	resp, err := http.Get(url)
	if err != nil {
		return []timepoint{}, err
	}

	return parseResponse(resp.Body)
}

func fetchFullDay(latitude float64, longitude float64) ([]timepoint, error) {
	url := fmt.Sprintf("https://graphdata.buienradar.nl/2.0/forecast/geo/rain24hour/?lat=%f&lon=%f", latitude, longitude)
	resp, err := http.Get(url)
	if err != nil {
		return []timepoint{}, err
	}

	return parseResponse(resp.Body)
}
