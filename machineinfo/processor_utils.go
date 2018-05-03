package machineinfo

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type CoreInfo struct {
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

func FormatLoadToMap() map[string][string] {
	var loadMap map[string]string = make(map[string]string)
	for i, load := range Loads {
		if i == 0 {
			loadMap["Total"] = fmt.Sprintf("%.2f", load)
		} else {
			loadMap[strconv.Itoa(i-1)] = fmt.Sprintf("%.2f", load)
		}
	}
	return loadMap
}

func CoreLoad(prevCoreInfo []CoreInfo, curCoreInfo []CoreInfo) []float64 {
	var loads []float64
	for i := 0; i < len(prevCoreInfo); i++ {
		coreLoad := loadAlgo(prevCoreInfo[i], curCoreInfo[i])
		loads = append(loads, coreLoad)
	}
	return loads
}


func loadAlgo(prevCoreInfo CoreInfo, curCoreInfo CoreInfo) float64 {
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

func ParseProcInfo() []CoreInfo {
	var cores []CoreInfo

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

func parseCpuLine(line []string) CoreInfo {
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

func fillCoreInfo(intList []int) CoreInfo {
	var CoreInfo CoreInfo = CoreInfo{
		intList[0], intList[1], intList[2], intList[3],
		intList[4], intList[5], intList[6], intList[7],
		intList[8], intList[9], intList[10]}

	return CoreInfo
}
