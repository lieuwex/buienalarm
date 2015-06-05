package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

type WeatherInfo struct {
	Start  float64   `json:"start"`
	Precip []float64 `json:"precip"`
}

func (info WeatherInfo) GetAmount() [25]float64 {
	var result [25]float64
	for i, val := range info.Precip {
		result[i] = math.Pow(10, float64((val-109)/32))
	}
	return result
}

func main() {
	resp, err := http.Get("http://www.buienalarm.nl/app/forecast.php?type=json&x=315&y=426")
	if err != nil {
		fmt.Println("Couldn't retreive forcast data from Buienalarm. ")
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var data WeatherInfo
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Couldn't parse JSON. ")
		panic(err)
	}

	t := time.Unix(int64(data.Start), 0)
	for _, val := range data.GetAmount() {
		fmt.Printf("%02d:%02d: %.1fmm/u\n", t.Hour(), t.Minute(), val)
		t = t.Add(5 * time.Minute)
	}
}
