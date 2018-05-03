package machineinfo

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const meminfoPath string = "/proc/meminfo"

type Memory struct {
	Total     int
	Used      int
	Free      int
	SwapTotal int
	SwapFree  int
	SwapUsed  int
}

func MemAllocation() Memory {
	var mem Memory = Memory{}

	m := parseMeminfo()
	mem.Total = m["MemTotal"]
	mem.Free = m["MemAvailable"]
	mem.Used = m["MemTotal"] - m["MemAvailable"]
	mem.SwapTotal = m["SwapTotal"]
	mem.SwapFree = m["SwapFree"]
	mem.SwapUsed = m["SwapTotal"] - m["SwapFree"]

	return mem
}

func parseMeminfo() map[string]int {
	var m map[string]int = make(map[string]int)

	meminfoFile, err := os.Open(meminfoPath)
	if err != nil {
		log.Fatal(err)
	}
	defer meminfoFile.Close()

	scanner := bufio.NewScanner(meminfoFile)
	for scanner.Scan() {
		var temp int

		line := strings.Fields(scanner.Text())
		temp, err = strconv.Atoi(line[1])
		if err != nil {
			log.Fatal(err)
		}

		line[0] = strings.Trim(line[0], ": ")
		m[line[0]] = temp
	}
	return m
}

func (mem Memory) FormatToMap() map[string]string {
	var memList map[string]string = make(map[string]string, 6)

	memList["Total"] = ConvertKilobytes(float64(mem.Total))
	memList["Used"] = ConvertKilobytes(float64(mem.Used))
	memList["Free"] = ConvertKilobytes(float64(mem.Free))
	memList["SwapTotal"] = ConvertKilobytes(float64(mem.SwapTotal))
	memList["SwapFree"] = ConvertKilobytes(float64(mem.SwapFree))
	memList["SwapUsed"] = ConvertKilobytes(float64(mem.SwapUsed))

	return memList
}

func ConvertKilobytes(kb float64) string {
	var size float64
	var sizeName string
	if kb < 1024 {
		size = kb
		sizeName = "KB"
	} else if kb > 1024 && kb < math.Pow(1024, 2) {
		size = kb / 1024
		sizeName = "MB"
	} else {
		size = kb / math.Pow(1024, 2)
		sizeName = "GB"
	}
	return fmt.Sprintf("%.2f %s", size, sizeName)
}
