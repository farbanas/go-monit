package main

import "fmt"

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
