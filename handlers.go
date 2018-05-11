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
	var loadsBars []string
	for i, _ := range loadsPercentage {
		loadsBars = append(loadsBars, DisplayPercentageBar(int(loadsPercentage[i]*100)))
	}
	mem := machineinfo.MemAllocation()
	memBar := DisplayPercentageBar(int((float64(mem.Used) / float64(mem.Total)) * 100))

	t, err := template.ParseFiles("templates/overview.html")
	if err != nil {
		log.Fatal(err)
	}

	t.Execute(w, struct {
		TotalLoad string
		Loads     []string
		Mem       string
	}{loadsBars[0], loadsBars[1:], memBar})
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
	memBar := DisplayPercentageBar(int((float64(mem.Used) / float64(mem.Total)) * 100))

	slackMsg["attachments"] = make([]map[string]string, 2)
	slackMsg["attachments"][0] = map[string]string{"text": fmt.Sprintf("MEM: %s", memBar)}
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
