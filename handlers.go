package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/farbanas/go-monit/machineinfo"
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

	w.Header().Set("Content-Type", "application/json")

	slackMsg := make(map[string][]map[string]string)

	/*
		r.ParseForm()
		command := r.Form.Get("command")
		text := r.Form.Get("text")
	*/
	mem := machineinfo.MemAllocation()
	memPercentage := int(mem.Used / mem.Total)

	slackMsg["attachments"] = make([]map[string]string, 2)
	slackMsg["attachments"][0] = map[string]string{"text": fmt.Sprintf("MEM: %s", DisplayPercentageBar(memPercentage))}
	slackMsg["attachments"][1] = map[string]string{"text": fmt.Sprintf("CPU: %s", DisplayPercentageBar(int(Loads[0]*100)))}

	json.NewEncoder(w).Encode(slackMsg)
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

	loadMap := FormatLoadToMap()
	json.NewEncoder(w).Encode(loadMap)
}
