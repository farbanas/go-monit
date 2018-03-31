package machineinfo

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
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

const procinfoPath string = "/proc/stat"

func NumCores() {}

func coreLoad() float64 {
	return 0.0
}

func ParseProcinfo() []coreInfo {
	var cores []coreInfo

	f, err := os.Open(procinfoPath)
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
