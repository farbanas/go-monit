package main

import (
	"html/template"
	"log"
	"net/http"

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
	}{totalLoad, loadsPercentage, mem.MemFormat()})
}

func slackMonitorHandler() {}

func memoryUsageHandler() {}

func loadSummaryHandler() {

}
