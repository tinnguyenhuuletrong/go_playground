package utils

import (
	"fmt"
	"log"
	"time"
)

func GetCurrentTime() int64 {
	return time.Now().Unix()
}

func HumanizeDurationAsString(from time.Time, to time.Time) string {
	var msSince = to.UnixMilli() - from.UnixMilli()
	log.Println("msSince: ", msSince)
	if msSince < 1000 {
		return fmt.Sprintf("%d ms", msSince)
	} else {
		var secSince = msSince / 1000
		log.Println("secSince ", secSince)
		if secSince < 60 {
			return fmt.Sprintf("%d sec", secSince)
		} else {
			var minSince = secSince / 60
			log.Println("minSince ", minSince)
			if minSince < 60 {
				return fmt.Sprintf("%d min", minSince)
			} else {
				var hourSince = minSince / 60
				log.Println("minSince ", hourSince)
				if hourSince < 60 {
					return fmt.Sprintf("%d hour", hourSince)
				} else {
					var daySince = hourSince / 24
					log.Println("daySince ", hourSince)
					return fmt.Sprintf("%d day", daySince)
				}
			}
		}
	}
}
