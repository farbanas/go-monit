package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"git.vingd.com/v-lab/go-monit/machineinfo"
)

func overviewHandler(w http.ResponseWriter, r *http.Request) {
	loadsPercentage := Loads
	for i, _ := range loadsPercentage {
		loadsPercentage[i] *= 100
	}
	totalLoad := loadsPercentage[0]
	loadsPercentage = loadsPercentage[1:]
	t, err := template.ParseFiles("templates/overview.html")
	if err != nil {
		log.Fatal(err)
	}
	mem := machineinfo.MemAllocation()

	t.Execute(w, struct {
		TotalLoad float64
		Loads     []float64
		Mem       map[string]string
	}{totalLoad, loadsPercentage, mem.FormatToMap()})
}

func slackMonitorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "This endpoint does not accept methods other than POST!")
	}
	defer r.Body.Close()

	var jsonMap map[string]interface{}

	fmt.Println(json.NewDecoder(r.Body).Decode(jsonMap))
}

func memoryUsageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		fmt.Fprintf(w, "This endpoint does not accept methods other than GET!")
	}
	mem := machineinfo.MemAllocation()
	json.NewEncoder(w).Encode(mem.FormatToMap())
}

func loadSummaryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		fmt.Fprintf(w, "This endpoint does not accept methods other than GET!")
	}

	var loadMap map[string]string = make(map[string]string)
	for i, load := range Loads {
		if i == 0 {
			loadMap["Total"] = fmt.Sprintf("%.2f", load)
		} else {
			loadMap[strconv.Itoa(i)] = fmt.Sprintf("%.2f", load)
		}
	}
	json.NewEncoder(w).Encode(loadMap)
}
