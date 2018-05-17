package utils

import (
	"fmt"
	"strconv"
)

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
