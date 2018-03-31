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

type memory struct {
	total     int
	used      int
	free      int
	swapTotal int
	swapFree  int
	swapUsed  int
}

func MemAllocation() memory {
	var mem memory = memory{}

	m := parseMeminfo()
	mem.total = m["MemTotal"]
	mem.free = m["MemFree"]
	mem.used = m["MemTotal"] - m["MemFree"]
	mem.swapTotal = m["SwapTotal"]
	mem.swapFree = m["SwapFree"]
	mem.swapUsed = m["SwapTotal"] - m["SwapFree"]

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

func ConvertKilobytes(kb float64) string {
	var size float64
	var sizeName string
	if kb < 1024 {
		size = kb
		sizeName = "kB"
	} else if kb > 1024 && kb < math.Pow(1024, 2) {
		size = kb / 1024
		sizeName = "mB"
	} else {
		size = kb / math.Pow(1024, 2)
		sizeName = "gB"
	}
	return fmt.Sprintf("%.2f %s", size, sizeName)
}
