package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/farbanas/go-monit/app/utils/machineinfo"
	"github.com/farbanas/go-monit/app/utils"
	"strings"
	"github.com/julienschmidt/httprouter"
)

func Overview(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	loadsPercentage := <- utils.Info.LoadChan
	var loadsBars []string
	for i := range loadsPercentage {
		loadsBars = append(loadsBars, utils.DisplayPercentageBar(int(loadsPercentage[i]*100)))
	}
	mem := machineinfo.MemAllocation()
	memBar := utils.DisplayPercentageBar(int((float64(mem.Used) / float64(mem.Total)) * 100))

	t, err := template.ParseFiles("app/views/overview.html")
	if err != nil {
		log.Fatal(err)
	}

	t.Execute(w, struct {
		TotalLoad string
		Loads     []string
		Mem       string
	}{loadsBars[0], loadsBars[1:], memBar})
}

func SlackMonitor(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "This endpoint does not accept methods other than POST!")
	}
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	slackMsg := make(map[string][]map[string]string)


	r.ParseForm()
	text := r.Form.Get("text")
	slackMsg["attachments"] = make([]map[string]string, 2)

	if strings.HasPrefix(text, "memory"){
		mem := machineinfo.MemAllocation()
		memBar := utils.DisplayPercentageBar(int((float64(mem.Used) / float64(mem.Total)) * 100))
		slackMsg["attachments"][0] = map[string]string{"text": fmt.Sprintf("MEM: %s", memBar)}
	} else if strings.HasPrefix(text,"load") {
		loads := <- utils.Info.LoadChan
		slackMsg["attachments"][0] = map[string]string{"text": fmt.Sprintf("CPU: %s", utils.DisplayPercentageBar(int(loads[0]*100)))}
	} else {
		mem := machineinfo.MemAllocation()
		memBar := utils.DisplayPercentageBar(int((float64(mem.Used) / float64(mem.Total)) * 100))
		loads := <- utils.Info.LoadChan
		slackMsg["attachments"][0] = map[string]string{"text": fmt.Sprintf("MEM: %s", memBar)}
		slackMsg["attachments"][1] = map[string]string{"text": fmt.Sprintf("CPU: %s", utils.DisplayPercentageBar(int(loads[0]*100)))}
	}

	json.NewEncoder(w).Encode(slackMsg)
}

func MemoryUsage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if r.Method != "GET" {
		fmt.Fprintf(w, "This endpoint does not accept methods other than GET!")
	}
	mem := machineinfo.MemAllocation()
	json.NewEncoder(w).Encode(mem.FormatToMap())
}

func LoadSummary(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if r.Method != "GET" {
		fmt.Fprintf(w, "This endpoint does not accept methods other than GET!")
	}

	loads := <- utils.Info.LoadChan
	loadMap := utils.FormatLoadToMap(loads)
	json.NewEncoder(w).Encode(loadMap)
}

func ProcessStatus(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if r.Method != "GET" {
		fmt.Fprintf(w, "This endpoint does not accept methods other than GET!")
	}

	processName := p.ByName("process")
	fmt.Println(processName)
}
