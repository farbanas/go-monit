package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/farbanas/go-monit/app/utils/machineinfo"
	"github.com/farbanas/go-monit/app/utils"
)

func OverviewHandler(w http.ResponseWriter, r *http.Request, loads []float64) {
	loadsPercentage := loads
	var loadsBars []string
	for i := range loadsPercentage {
		loadsBars = append(loadsBars, utils.DisplayPercentageBar(int(loadsPercentage[i]*100)))
	}
	mem := machineinfo.MemAllocation()
	memBar := utils.DisplayPercentageBar(int((float64(mem.Used) / float64(mem.Total)) * 100))

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

func SlackMonitorHandler(w http.ResponseWriter, r *http.Request, loads []float64) {
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
	memBar := utils.DisplayPercentageBar(int((float64(mem.Used) / float64(mem.Total)) * 100))

	slackMsg["attachments"] = make([]map[string]string, 2)
	slackMsg["attachments"][0] = map[string]string{"text": fmt.Sprintf("MEM: %s", memBar)}
	slackMsg["attachments"][1] = map[string]string{"text": fmt.Sprintf("CPU: %s", utils.DisplayPercentageBar(int(loads[0]*100)))}

	json.NewEncoder(w).Encode(slackMsg)
}

func MemoryUsageHandler(w http.ResponseWriter, r *http.Request, loads []float64) {
	if r.Method != "GET" {
		fmt.Fprintf(w, "This endpoint does not accept methods other than GET!")
	}
	mem := machineinfo.MemAllocation()
	json.NewEncoder(w).Encode(mem.FormatToMap())
}

func LoadSummaryHandler(w http.ResponseWriter, r *http.Request, loads []float64) {
	if r.Method != "GET" {
		fmt.Fprintf(w, "This endpoint does not accept methods other than GET!")
	}

	loadMap := utils.FormatLoadToMap(loads)
	json.NewEncoder(w).Encode(loadMap)
}
