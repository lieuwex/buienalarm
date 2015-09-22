package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
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
	resp, err := http.Get("http://www.buienalarm.nl/location/Wassenaar")
	if err != nil {
		fmt.Println("Couldn't retreive forcast data from Buienalarm. ")
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	for _, line := range strings.Split(string(body), "\n") {
		if strings.Contains(line, "locationdata['forecast']") {
			x := strings.Split(line, " ")
			raw := strings.Replace(x[len(x)-1], ";", "", -1)

			var data WeatherInfo
			err = json.Unmarshal([]byte(raw), &data)
			if err != nil {
				fmt.Println("Couldn't parse JSON. ")
				panic(err)
			}

			t := time.Unix(int64(data.Start), 0)
			for _, val := range data.GetAmount() {
				fmt.Printf("%02d:%02d: %.1fmm/u\n", t.Hour(), t.Minute(), val)
				t = t.Add(5 * time.Minute)
			}

			break
		}
	}

}