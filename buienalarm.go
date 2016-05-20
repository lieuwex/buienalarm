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
	Start  float64   `json:"start"`
	Precip []float64 `json:"precip"`
}

func (info weatherInfo) getAmount() [25]float64 {
	var result [25]float64
	for i, val := range info.Precip {
		result[i] = math.Pow(10, float64((val-109)/32))
	}
	return result
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <location>\n", strings.Join(os.Args, " "))
		os.Exit(1)
	}
	location := os.Args[1]

	resp, err := http.Get("http://www.buienalarm.nl/location/" + location)
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

			var data weatherInfo
			err = json.Unmarshal([]byte(raw), &data)
			if err != nil {
				fmt.Println("Couldn't parse JSON. ")
				panic(err)
			}

			t := time.Unix(int64(data.Start), 0)
			for _, val := range data.getAmount() {
				fmt.Printf("%02d:%02d: %.1fmm/u\n", t.Hour(), t.Minute(), val)
				t = t.Add(5 * time.Minute)
			}

			break
		}
	}

}
