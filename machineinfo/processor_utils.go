package machineinfo

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type coreInfo struct {
	coreNum    int
	user       int
	nice       int
	system     int
	idle       int
	iowait     int
	irq        int
	softirq    int
	steal      int
	guest      int
	guest_nice int
}

const procInfoPath string = "/proc/stat"

func CoreLoad(loadChan chan []float64) {
	var prevCoreInfo []coreInfo
	var curCoreInfo []coreInfo

	curCoreInfo = ParseProcInfo()
	time.Sleep(1 * time.Second)
	for {
		prevCoreInfo = curCoreInfo
		curCoreInfo = ParseProcInfo()
		go func() {
			calcCoreLoad(loadChan, prevCoreInfo, curCoreInfo)
		}()
		time.Sleep(1 * time.Second)
	}
}

func calcCoreLoad(loadChan chan []float64, prevCoreInfo []coreInfo, curCoreInfo []coreInfo) {
	var loads []float64
	for i := 0; i < len(prevCoreInfo); i++ {
		coreLoad := loadAlgo(prevCoreInfo[i], curCoreInfo[i])
		loads = append(loads, coreLoad)
	}
	loadChan <- loads
}

func loadAlgo(prevCoreInfo coreInfo, curCoreInfo coreInfo) float64 {
	prevIdle := prevCoreInfo.idle + prevCoreInfo.iowait
	idle := curCoreInfo.idle + curCoreInfo.iowait

	prevNonIdle := prevCoreInfo.user + prevCoreInfo.nice + prevCoreInfo.system + prevCoreInfo.irq + prevCoreInfo.softirq + prevCoreInfo.steal
	nonIdle := curCoreInfo.user + curCoreInfo.nice + curCoreInfo.system + curCoreInfo.irq + curCoreInfo.softirq + curCoreInfo.steal

	prevTotal := prevIdle + prevNonIdle
	total := idle + nonIdle

	totald := total - prevTotal
	idled := idle - prevIdle

	load := float64(totald-idled) / float64(totald)
	return load
}

func ParseProcInfo() []coreInfo {
	var cores []coreInfo

	f, err := os.Open(procInfoPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if strings.HasPrefix(line[0], "cpu") {
			cores = append(cores, parseCpuLine(line))
		}
	}
	return cores
}

func parseCpuLine(line []string) coreInfo {
	var intList []int

	line[0] = strings.TrimLeft(line[0], "cpu ")
	for _, el := range line {
		var err error
		var intEl int

		if el != "" {
			intEl, err = strconv.Atoi(el)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			intEl = -1
		}
		intList = append(intList, intEl)
	}
	return fillCoreInfo(intList)
}

func fillCoreInfo(intList []int) coreInfo {
	var coreInfo coreInfo = coreInfo{
		intList[0], intList[1], intList[2], intList[3],
		intList[4], intList[5], intList[6], intList[7],
		intList[8], intList[9], intList[10]}

	return coreInfo
}
