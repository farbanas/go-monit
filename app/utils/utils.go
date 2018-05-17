package utils

import (
	"fmt"
	"strconv"
	"github.com/farbanas/go-monit/app/utils/machineinfo"
	"time"
)

type ComputerInfo struct {
	LoadChan chan []float64
}

var Info = ComputerInfo{make(chan []float64, 5)}

func DisplayPercentageBar(percentage int) string {
	poundsNum := int(percentage / 5)
	pounds := ""
	for i := 0; i < 20; i++ {
		if i < poundsNum {
			pounds += "#"
		} else {
			pounds += "."
		}
	}

	percentageBar := fmt.Sprintf("[%s] %d%%", pounds, percentage)
	return percentageBar
}

func FormatLoadToMap(loads []float64) map[string]string {
	var loadMap = make(map[string]string)
	for i, load := range loads {
		if i == 0 {
			loadMap["Total"] = fmt.Sprintf("%.2f%%", load*100)
		} else {
			loadMap[strconv.Itoa(i-1)] = fmt.Sprintf("%.2f%%", load*100)
		}
	}
	return loadMap
}

func (info ComputerInfo) CoreLoadFeed() {
	var prevCoreInfo []machineinfo.CoreInfo
	var curCoreInfo []machineinfo.CoreInfo

	curCoreInfo = machineinfo.ParseProcInfo()
	time.Sleep(1 * time.Second)
	for range time.Tick(1 * time.Second) {
		prevCoreInfo = curCoreInfo
		curCoreInfo = machineinfo.ParseProcInfo()
		info.LoadChan <- machineinfo.CoreLoad(prevCoreInfo, curCoreInfo)
	}
}

func Init() {
	go Info.CoreLoadFeed()
}
