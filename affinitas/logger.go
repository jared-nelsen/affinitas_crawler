package affinitas

import "time"

func getTimeStr() string {
	currentTime := time.Now()
	return currentTime.Format("15:04:05")
}

func Log(message string) {
	println(" - " + getTimeStr() + " -- " + message)
}
