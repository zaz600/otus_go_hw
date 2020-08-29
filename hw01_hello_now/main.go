package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

const ntpServer = "0.beevik-ntp.pool.ntp.org"
const timeFormat = "2006-01-02 15:04:05 -0700 MST"

func main() {
	ntpTime, err := ntp.Time(ntpServer)
	if err != nil {
		log.Fatalf("Error get ntp time %s\n", err)
	}
	localTime := time.Now()

	fmt.Printf("current time: %s\n", localTime.Format(timeFormat))
	fmt.Printf("exact time: %s\n", ntpTime.Format(timeFormat))
}
