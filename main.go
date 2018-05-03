package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/farbanas/go-monit/machineinfo"
)

var Loads []float64

func main() {
	loadChan := make(chan []float64)

	go coreLoadFeed(loadChan)
	go collectLoad(loadChan)

	http.HandleFunc("/", overviewHandler)
	http.HandleFunc("/webhooks/slack/monitor", slackMonitorHandler)
	http.HandleFunc("/webhooks/load", loadSummaryHandler)
	http.HandleFunc("/webhooks/memory", memoryUsageHandler)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func coreLoadFeed(loadChan chan []float64) {
	var prevCoreInfo []machineinfo.CoreInfo
	var curCoreInfo []machineinfo.CoreInfo

	curCoreInfo = machineinfo.ParseProcInfo()
	time.Sleep(1 * time.Second)
	for _ = range time.Tick(1 * time.Second) {
		go func() {
			prevCoreInfo = curCoreInfo
			curCoreInfo = machineinfo.ParseProcInfo()
			loadChan <- machineinfo.CoreLoad(prevCoreInfo, curCoreInfo)
		}()
	}
}

func collectLoad(loadChan chan []float64) {
	for {
		Loads = <-loadChan
	}
}

func FormatLoadToMap() map[string]string {
	var loadMap map[string]string = make(map[string]string)
	for i, load := range Loads {
		if i == 0 {
			loadMap["Total"] = fmt.Sprintf("%.2f%%", load*100)
		} else {
			loadMap[strconv.Itoa(i-1)] = fmt.Sprintf("%.2f%%", load*100)
		}
	}
	return loadMap
}
